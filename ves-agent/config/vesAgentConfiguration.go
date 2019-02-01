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
	"github.com/nokia/onap-vespa/govel"
)

// VESAgentConfiguration parameters
type VESAgentConfiguration struct {
	PrimaryCollector govel.CollectorConfiguration    `mapstructure:"primaryCollector"`
	BackupCollector  govel.CollectorConfiguration    `mapstructure:"backupCollector,omitempty"`
	Heartbeat        HeartbeatConfiguration    `mapstructure:"heartbeat,omitempty"`
	Measurement      MeasurementConfiguration  `mapstructure:"measurement,omitempty"`
	Event            govel.EventConfiguration        `mapstructure:"event,omitempty"`
	AlertManager     AlertManagerConfiguration `mapstructure:"alertManager,omitempty"`
	Cluster          *ClusterConfiguration     `mapstructure:"cluster"` // Optional cluster config. If absent, fallbacks to single node mode
	Debug            bool                      `mapstructure:"debug,omitempty"`
	CaCert           string                    `mapstructure:"caCert,omitempty"` // Root certificate content
	DataDir          string                    `mapsctructure:"datadir"`         // Path to directory containing data
}
