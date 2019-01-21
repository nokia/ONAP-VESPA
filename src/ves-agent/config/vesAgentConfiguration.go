package config

// VESAgentConfiguration parameters
type VESAgentConfiguration struct {
	PrimaryCollector CollectorConfiguration    `mapstructure:"primaryCollector"`
	BackupCollector  CollectorConfiguration    `mapstructure:"backupCollector,omitempty"`
	Heartbeat        HeartbeatConfiguration    `mapstructure:"heartbeat,omitempty"`
	Measurement      MeasurementConfiguration  `mapstructure:"measurement,omitempty"`
	Event            EventConfiguration        `mapstructure:"event,omitempty"`
	AlertManager     AlertManagerConfiguration `mapstructure:"alertManager,omitempty"`
	Cluster          *ClusterConfiguration     `mapstructure:"cluster"` // Optional cluster config. If absent, fallbacks to single node mode
	Debug            bool                      `mapstructure:"debug,omitempty"`
	CaCert           string                    `mapstructure:"caCert,omitempty"` // Root certificate content
	DataDir          string                    `mapsctructure:"datadir"`         // Path to directory containing data
}
