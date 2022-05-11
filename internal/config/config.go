package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"in-backend/pkg/logging"
	"sync"
)

type Config struct {
	IsDebug  *bool    `yaml:"is_debug"`
	Postgres Postgres `yaml:"postgres"`
	Listen   Listen   `yaml:"listen"`
}

type Listen struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"pass"`
	DB       string `yaml:"db"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
