package main

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	MigrationsPath string `env:"MIGRATIONS_PATH" env-default:"./migrations"`
	DbPath         string `env:"DB_PATH" env-default:"./data/arcane.db"`
	RedisHost      string `env:"REDIS_HOST" env-default:"127.0.0.1"`
	RedisPort      string `env:"REDIS_PORT" env-default:"6379"`
	RedisPassword  string `env:"REDIS_PASSWORD" env-default:"abracadabra"`
	Port           string `env:"PORT" env-default:"5000"`
	LogLevel       string `env:"LOG_LEVEL" env-default:"info"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	return &cfg, nil
}
