package config

// CollectorConfiguration parameters
type CollectorConfiguration struct {
	FQDN       string `mapstructure:"fqdn"`
	Port       int    `mapstructure:"port"`
	Topic      string `mapstructure:"topic"`
	User       string `mapstructure:"user"`
	Password   string `mapstructure:"password"`
	PassPhrase string `mapsctructure:"passphrase,omitempty"` // passPhrase used to encrypt collector password in file
}
