package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const configPath = ".env"

type Config struct {
	ApiKey      string        `env:"APIKEY"`
	BotKey      string        `env:"BOTKEY"`
	TimeOut     time.Duration `env:"TIMEOUT"`
	Database    Database
	RedisConfig RedisConfig
}

type Database struct {
	Host     string `env:"DATABASE_HOST"`
	Port     int    `env:"DATABASE_PORT"`
	User     string `env:"DATABASE_USER"`
	Password string `env:"DATABASE_PASSWORD"`
	DBName   string `env:"DATABASE_DBNAME"`
	DBCon    string `env:"DB_CON"`
}

type RedisConfig struct {
	Password string        `env:"REDIS_PASSWORD"`
	Host     string        `env:"REDIS_HOST"`
	Port     string        `env:"REDIS_PORT"`
	DB       int           `env:"REDIS_DB"`
	TTL      time.Duration `env:"REDIS_TTL"`
}

func Load() *Config {
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Cannot read сonfig: %s", err)
	}
	return &cfg
}
