package config

import (
	"sync"
	"time"
)

var instance *Config
var once sync.Once

type Config struct {
	Port              string
	PruneIntervalHour time.Duration
}

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{
			Port:              "8081",
			PruneIntervalHour: time.Hour * 1,
		}
	})

	return instance
}
