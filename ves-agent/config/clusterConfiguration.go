package config

import (
	"net"

	"github.com/hashicorp/raft"
)

// Peer is the configuration of a single raft peer
type Peer struct {
	ID      string `mapstructure:"id"`      // Peer ID
	Address string `mapstructure:"address"` // Peer address and port
}

// Server converts the peer configuration into a raft server
func (pcfg Peer) Server() raft.Server {
	return raft.Server{ID: raft.ServerID(pcfg.ID), Address: raft.ServerAddress(pcfg.Address)}
}

// TCPAddr resolves the peer address &nd port into a TCP address
func (pcfg *Peer) TCPAddr() (*net.TCPAddr, error) {
	return net.ResolveTCPAddr("tcp", pcfg.Address)
}

// Peers is a list of peer configuration
type Peers []Peer

// Servers converts the list of peer configuration
// into a list of raft servers
func (cfg Peers) Servers() []raft.Server {
	servers := make([]raft.Server, len(cfg))
	for i, peer := range cfg {
		servers[i] = peer.Server()
	}
	return servers
}

// GetPeer find and returns the peer configuration
// for a given peer ID
func (cfg Peers) GetPeer(id string) (Peer, bool) {
	for _, p := range cfg {
		if p.ID == id {
			return p, true
		}
	}
	return Peer{}, false
}

// ClusterConfiguration is the configuration of the
// raft cluster
type ClusterConfiguration struct {
	ID          string `mapstructure:"id"`          // Local node ID
	Peers       Peers  `mapstructure:"peers"`       // List of cluster's node
	Debug       bool   `mapstructure:"debug"`       // Display raft log messages
	DisplayLogs bool   `mapstructure:"displayLogs"` // Display replication log entries
}
