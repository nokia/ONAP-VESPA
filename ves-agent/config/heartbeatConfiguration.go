package config

import (
	"time"
)

// HeartbeatConfiguration parameters
type HeartbeatConfiguration struct {
	DefaultInterval time.Duration `mapstructure:"defaultInterval,omitempty"`
}
