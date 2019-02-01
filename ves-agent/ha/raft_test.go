/*
	Copyright 2019 Nokia

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

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
