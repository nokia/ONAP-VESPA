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
