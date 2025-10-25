package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port string `yaml:"port" env:"PORT" env-default:"5432"`
}

func InitCfg() *Config {
	cfg := Config{}

	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		log.Fatalf("read config error: %s", err.Error())
	}

	return &cfg
}
