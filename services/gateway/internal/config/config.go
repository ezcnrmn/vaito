package config

type Config struct {
	Port   string
	Rabbit struct {
		User string
		Pass string
		URL  string
	}
}
