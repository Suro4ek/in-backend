package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"in-backend/pkg/logging"
	"sync"
)

type Config struct {
	IsDebug  *bool    `yaml:"is_debug" env:"ISDEBUG"`
	Postgres Postgres `yaml:"postgres"`
	Secret   string   `yaml:"secret" env:"JWT_SECRET"`
	Pattern  Patterns `yaml:"patterns"`
}

type Patterns struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Postgres struct {
	Host     string `yaml:"host" env:"POSTGRES_HOST"`
	Port     string `yaml:"port" env:"POSTGRES_PORT"`
	User     string `yaml:"user" env:"POSTGRES_USER"`
	Password string `yaml:"pass" env:"POSTGRES_PASS"`
	DB       string `yaml:"db" env:"POSTGRES_DB"`
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
		if err := cleanenv.ReadEnv(instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
