package config

type ConsulConfig struct {
	Enabled bool   `envconfig:"ENABLED" default:"false"`
	Addr    string `envconfig:"ADDR" default:"0.0.0.0:8500"`
}

type Config struct {
	AppAddr       string       `envconfig:"APP_ADDR" default:":8083"`
	MetadataAddrs []string     `envconfig:"METADATA_ADDRS" default:"localhost:8081"`
	RatingAddrs   []string     `envconfig:"RATING_ADDRS" default:"localhost:8082"`
	Consul        ConsulConfig `envconfig:"CONSUL"`
}
