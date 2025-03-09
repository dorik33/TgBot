package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPath = "config.yaml"

type Config struct {
	BotKey   string        `yaml:"botkey"`
	TimeOut  time.Duration `yaml:"timeout"`
	Database Database      `yaml:"database"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

func Load() *Config {
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Cannot read konfig: %s", err)
	}
	return &cfg
}
