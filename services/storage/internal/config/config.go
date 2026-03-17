package config

import "time"

type Config struct {
	Port       string
	GatewayUrl string
	DB         struct {
		DSN          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  time.Duration
	}
}
