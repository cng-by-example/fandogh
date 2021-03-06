package config

type Config struct {
	Trace `koanf:"trace"`
}

type Trace struct {
	Enabled bool `koanf:"enabled"`
	Agent   `koanf:"agent"`
}

type Agent struct {
	Host string `koanf:"host"`
	Port string `koanf:"port"`
}
