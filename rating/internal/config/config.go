package config

type ConsulConfig struct {
	Enabled bool   `envconfig:"ENABLED" default:"false"`
	Addr    string `envconfig:"ADDR" default:"0.0.0.0:8500"`
}

type Config struct {
	AppAddr string       `envconfig:"APP_ADDR" default:":8082"`
	Consul  ConsulConfig `envconfig:"CONSUL"`
}
