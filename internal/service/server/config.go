package server

import (
	"flag"
	"github.com/kyrare/ya-metrics/internal/domain/utils"
)

type Config struct {
	address string
}

func LoadConfig() (Config, error) {
	addr := utils.GetParameter("a", "ADDRESS", "0.0.0.0:8080", "Адрес сервера (по умолчанию 0.0.0.0:8080)")

	flag.Parse()

	return Config{
		address: *addr,
	}, nil
}
