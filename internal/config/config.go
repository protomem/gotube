package config

type HTTP struct {
	Host string
	Port int
}

type Config struct {
	HTTP
}

func New() Config {
	return Config{
		HTTP: HTTP{
			Host: "0.0.0.0",
			Port: 8080,
		},
	}
}
