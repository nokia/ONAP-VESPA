package ha

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	"ves-agent/config"

	"github.com/hashicorp/raft-boltdb"

	"github.com/hashicorp/raft"
	log "github.com/sirupsen/logrus"
)

func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

// Cluster is a replicated AgentState in a raft cluster.
// It provides consistent write operations across the cluster, and those
// write operations are permitted only on the leader node.
type Cluster struct {
	id       string
	raft     *raft.Raft
	leaderCh <-chan bool
	fsm      *FSM
}

// NewCluster creates and start a new cluster around `state`.
// Cluster topology is included in `cfg`, but if `cfg` is null, the cluster will fallsback into single-node mode.
// On first creation, the cluster is boostrapped
func NewCluster(datadir string, cfg *config.ClusterConfiguration, state SnapshotableAgentState) (*Cluster, error) {
	leaderCh := make(chan bool, 1024)
	fsm := NewFSM(state, cfg != nil && cfg.DisplayLogs)

	conf := raft.DefaultConfig()
	// conf.SnapshotInterval = 10 * time.Second
	// conf.SnapshotThreshold = 100
	// conf.TrailingLogs = 10
	var logOutput io.Writer
	if cfg != nil && cfg.Debug {
		logOutput = log.StandardLogger().Writer()
	} else {
		logOutput = ioutil.Discard
	}
	conf.LogOutput = logOutput
	conf.NotifyCh = leaderCh
	var transport raft.Transport
	var servers []raft.Server
	var snapshotStore raft.SnapshotStore
	var store interface {
		raft.StableStore
		raft.LogStore
	}
	var needBootstrap bool

	if cfg != nil && len(cfg.Peers) > 0 {
		log.Info("Initializing Raft cluster")
		conf.LocalID = raft.ServerID(cfg.ID)
		servers = cfg.Peers.Servers()
		myself, ok := cfg.Peers.GetPeer(cfg.ID)
		if !ok {
			return nil, fmt.Errorf("Bad cluster configuration: No peer with id %s found", cfg.ID)
		}

		advertise, err := myself.TCPAddr()
		if err != nil {
			return nil, err
		}
		transport, err = raft.NewTCPTransport(myself.Address, advertise, 1000, 5*time.Second, logOutput)
		if err != nil {
			return nil, err
		}
		baseDir := filepath.Join(datadir, "raft", string(conf.LocalID))
		if err = os.MkdirAll(baseDir, 0750); err != nil {
			return nil, err
		}
		storeFile := filepath.Join(baseDir, "store.db")
		needBootstrap = !fileExists(storeFile)
		store, err = raftboltdb.NewBoltStore(storeFile)
		if err != nil {
			return nil, err
		}
		snapshotStore, err = raft.NewFileSnapshotStore(baseDir, 5, logOutput)
		if err != nil {
			return nil, err
		}
	} else {
		log.Warn("No cluster config or no peers found. Fallback to single node mode")
		conf.LocalID = raft.ServerID("single-node")
		var addr raft.ServerAddress
		addr, transport = raft.NewInmemTransport(raft.ServerAddress(""))
		servers = append(servers, raft.Server{ID: conf.LocalID, Address: addr})
		store = raft.NewInmemStore()
		snapshotStore = raft.NewInmemSnapshotStore()
		needBootstrap = true
	}

	node, err := raft.NewRaft(conf, fsm, store, store, snapshotStore, transport)
	if err != nil {
		return nil, err
	}

	log.Infof("Cluster %s created", node.String())
	cluster := &Cluster{
		id:       string(conf.LocalID),
		raft:     node,
		leaderCh: leaderCh,
		fsm:      fsm,
	}

	if needBootstrap {
		return cluster, cluster.bootstrap(servers)
	}
	return cluster, nil
}

// Stats return a map with stats about the Raft cluster.
// This should be used only for debugging purpose. Do net expect
// this interface to remain stable over time
func (cluster *Cluster) Stats() map[string]string {
	s := cluster.raft.Stats()
	s["leader"] = string(cluster.raft.Leader())
	return s
}

// LeaderCh returns a buffered channel which receive cluster leadership
// changes for the current node. It MUST be consummed
func (cluster *Cluster) LeaderCh() <-chan bool {
	return cluster.leaderCh
}

func (cluster *Cluster) bootstrap(peers []raft.Server) error {
	log.Info("Bootstrapping Raft Cluster")
	return cluster.raft.BootstrapCluster(raft.Configuration{
		Servers: peers,
	}).Error()
}

func (cluster *Cluster) apply(cmd StateCmd) (interface{}, error) {
	data, err := json.Marshal(&cmd)
	if err != nil {
		return nil, err
	}
	f := cluster.raft.Apply(data, 5*time.Second)
	if err := f.Error(); err != nil {
		return nil, err
	}
	switch resp := f.Response().(type) {
	case error:
		return nil, resp
	default:
		return resp, nil
	}
}

// Shutdown stops the current raft node, and all associated goroutines
func (cluster *Cluster) Shutdown() error {
	return cluster.raft.Shutdown().Error()
}

// NextMeasurementIndex return the next event index and increments it
func (cluster *Cluster) NextMeasurementIndex() (int64, error) {
	idx, err := cluster.apply(StateCmd{Type: IncrementMeasIdx})
	if err != nil {
		return 0, err
	}
	return idx.(int64), nil
}

// NextHeartbeatIndex return the next event index and increments it
func (cluster *Cluster) NextHeartbeatIndex() (int64, error) {
	idx, err := cluster.apply(StateCmd{Type: IncrementHeartbeatIdx})
	if err != nil {
		return 0, err
	}
	return idx.(int64), nil
}

// NextRun returns the time at which next execution should occure
func (cluster *Cluster) NextRun(sched string) time.Time {
	return cluster.fsm.NextRun(sched)
}

// UpdateNextRun set the time of the next execution
func (cluster *Cluster) UpdateNextRun(sched string, next time.Time) error {
	nxt := next.Unix()
	_, err := cluster.apply(StateCmd{Type: UpdateScheduler, UpdateScheduler: &UpdateSchedulerFields{Name: sched, Next: &nxt}})
	return err
}

// Interval returns the scheduler exceution interval
func (cluster *Cluster) Interval(sched string) time.Duration {
	return cluster.fsm.Interval(sched)
}

// UpdateInterval set a new execution interval for the scheduler
func (cluster *Cluster) UpdateInterval(sched string, interval time.Duration) error {
	_, err := cluster.apply(StateCmd{Type: UpdateScheduler, UpdateScheduler: &UpdateSchedulerFields{Name: sched, Interval: &interval}})
	return err
}

// UpdateScheduler set both interval and next execution time for the scheduler
func (cluster *Cluster) UpdateScheduler(sched string, interval time.Duration, next time.Time) error {
	nxt := next.Unix()
	_, err := cluster.apply(StateCmd{Type: UpdateScheduler, UpdateScheduler: &UpdateSchedulerFields{Name: sched, Interval: &interval, Next: &nxt}})
	return err
}

// GetFaultSn return the fault sequence number
func (cluster *Cluster) GetFaultSn(faultID int32) int64 {
	return cluster.fsm.GetFaultSn(faultID)
}

// GetFaultStartEpoch return the startEpoch
func (cluster *Cluster) GetFaultStartEpoch(faultID int32) int64 {
	return cluster.fsm.GetFaultStartEpoch(faultID)
}

// IncrementFaultSn increments the fault sequence number
func (cluster *Cluster) IncrementFaultSn(faultID int32) error {
	//fmt.Printf("raft IncrementFaultSn faultID:%d", faultID)
	sn := int64(1)
	//faultName := strconv.FormatInt(faultID(int64), 10)
	_, err := cluster.apply(StateCmd{Type: UpdateFault, UpdateFault: &UpdateFaultFields{FaultID: &faultID, SequenceNumber: &sn}})
	return err
}

// SetFaultStartEpoch set the startEpoch
func (cluster *Cluster) SetFaultStartEpoch(faultID int32, epoch int64) error {
	//faultName := strconv.Itoa(faultID)
	_, err := cluster.apply(StateCmd{Type: UpdateFault, UpdateFault: &UpdateFaultFields{FaultID: &faultID, StartEpoch: &epoch}})
	return err
}

// InitAlertInfos update the alertInfos map
func (cluster *Cluster) InitAlertInfos(faultID int32) error {
	return cluster.fsm.InitAlertInfos(faultID)
}

// NextFaultIndex return the Fault Index and increments it
func (cluster *Cluster) NextFaultIndex() (int32, error) {
	idx, err := cluster.apply(StateCmd{Type: IncrementFaultIdx})
	if err != nil {
		return 0, err
	}
	return idx.(int32), nil
}

// GetFaultInStorage checks if faultName already associated to an index
func (cluster *Cluster) GetFaultInStorage(faultName string) int32 {
	return cluster.fsm.GetFaultInStorage(faultName)
}

// StoreFaultInStorage stores the index associated to the faultName
func (cluster *Cluster) StoreFaultInStorage(faultName string, faultID int32) error {
	//fmt.Printf("raft msg StoreFaultInStorage faultName:%s faultID:%d", faultName, faultID)
	_, err := cluster.apply(StateCmd{Type: UpdateFault, UpdateFault: &UpdateFaultFields{FaultID: &faultID, FaultName: faultName}})
	return err
}

// DeleteFaultInStorage delete Fault in storage
func (cluster *Cluster) DeleteFaultInStorage(faultName string) error {
	//fmt.Printf("raft msg DeleteFaultInStorage faultName:%s ", faultName)
	_, err := cluster.apply(StateCmd{Type: DeleteFault, DeleteFault: &DeleteFaultFields{FaultName: faultName}})
	return err
}
