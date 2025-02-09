package config

import (
	"time"

	"github.com/Rolan335/Finop/internal/repository/postgres"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Port           string        `env:"PORT" envDefault:"8080"`
	RequestTimeout time.Duration `env:"REQUEST_TIMEOUT" envDefault:"10s"`
	GinMode        string        `env:"GIN_MODE" envDefault:"debug"`
	DB             postgres.Config
	Migration      postgres.MigrationConfig
}

func MustNewConfig(pathToEnv string) *Config {
	//load env file
	if err := godotenv.Load(pathToEnv); err != nil {
		panic("failed to load env file: " + err.Error())
	}

	//parse env file
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		panic("failed to parse env: " + err.Error())
	}
	return &cfg
}
