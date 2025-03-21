package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sunzhqr/sharedlibprompt/pkg/postgres"
)

type Config struct {
	Postgres postgres.Config `yaml:"POSTGRES" env:"POSTGRES"`
	GRPCPort int             `yaml:"GRPC_PORT" env:"GRPC_PORT" env-default:"50051"`
	RestPort int             `yaml:"REST_PORT" env:"REST_PORT" env-default:"8081`
}

func New() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig("./config/config.yaml", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
