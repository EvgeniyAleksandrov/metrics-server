package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
)

type Server struct {
	Address       string `env:"ADDRESS"`
	StoreInterval int    `env:"STORE_INTERVAL"`
	FilePath      string `env:"FILE_STORAGE_PATH"`
	Restore       bool   `env:"RESTORE"`
	DatabaseDSN   string `env:"DATABASE_DSN"`
}

func ParseServerConfig() (*Server, error) {
	serverConfig := &Server{}

	flag.StringVar(
		&serverConfig.Address,
		"a",
		"localhost:8080",
		"The address and port on which the server listens for connections.",
	)

	flag.IntVar(&serverConfig.StoreInterval, "i", 300, "intervale for save storage to disk")
	flag.StringVar(&serverConfig.FilePath, "f", "db.json", "storage file path")
	flag.BoolVar(&serverConfig.Restore, "r", true, "need to restore data from file")
	flag.StringVar(&serverConfig.DatabaseDSN, "d", "", "postgress database connection string")

	flag.Parse()

	if err := env.Parse(serverConfig); err != nil {
		return nil, fmt.Errorf("parse env variables: %w", err)
	}

	return serverConfig, nil
}
