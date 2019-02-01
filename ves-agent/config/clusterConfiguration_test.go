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

package config

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ClusterConfigTestSuite struct {
	suite.Suite
	cfg *ClusterConfiguration
}

func TestClusterConfig(t *testing.T) {
	suite.Run(t, new(ClusterConfigTestSuite))
}

func (s *ClusterConfigTestSuite) SetupTest() {
	s.cfg = &ClusterConfiguration{
		ID: "12",
		Peers: Peers{
			{ID: "1", Address: "127.0.0.1:6565"},
			{ID: "2", Address: "127.0.0.2:6565"},
			{ID: "3", Address: "127.0.0.3:6565"},
		},
	}
}

func (s *ClusterConfigTestSuite) TestGetPeer() {
	p, ok := s.cfg.Peers.GetPeer("3")
	s.True(ok)
	s.Equal("3", p.ID)

	p, ok = s.cfg.Peers.GetPeer("42")
	s.False(ok)
}

func (s *ClusterConfigTestSuite) TestServersConvert() {
	servers := s.cfg.Peers.Servers()
	s.NotNil(servers)
	s.Len(servers, 3)
}

func (s *ClusterConfigTestSuite) TestPeerTcpResolve() {
	addr, err := s.cfg.Peers[0].TCPAddr()
	s.NoError(err)
	s.NotNil(addr)
}
