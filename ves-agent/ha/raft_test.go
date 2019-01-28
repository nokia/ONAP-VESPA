package ha

import (
	"os"
	"testing"
	"github.com/nokia/onap-vespa/ves-agent/config"

	"github.com/stretchr/testify/suite"
)

func TestRaft(t *testing.T) {
	suite.Run(t, &StateTestSuite{
		setup: func(s *StateTestSuite) {
			var err error
			cfg := &config.ClusterConfiguration{Debug: true, DisplayLogs: true}
			s.state, err = NewCluster("./test_datadir", cfg, NewInMemState())
			if err != nil {
				panic(err)
			}
			for !<-s.state.(*Cluster).LeaderCh() {
			}
		},
		teardown: func(s *StateTestSuite) {
			s.NoError(s.state.(*Cluster).Shutdown())
			s.NoError(os.RemoveAll("./test_datadir"))
		},
	})
}
