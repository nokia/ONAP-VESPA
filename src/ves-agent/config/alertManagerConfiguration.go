package config

// AlertManagerConfiguration parameters
type AlertManagerConfiguration struct {
	Bind     string `yaml:"bind,omitempty"`
	Path     string `yaml:"path,omitempty"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
