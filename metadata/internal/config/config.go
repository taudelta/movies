package config

type Config struct {
	AppAddr string `envconfig:"APP_ADDR" default:":8081"`
}
