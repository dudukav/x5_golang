package config

import (
	"os"
	"strings"
)

const DefaultHTTPAddr = ":8080"

type Config struct {
	HTTPAddr string
}

func Load() Config {
	httpAddr := strings.TrimSpace(os.Getenv("HTTP_ADDR"))
	if httpAddr == "" {
		httpAddr = DefaultHTTPAddr
	}

	return Config{HTTPAddr: httpAddr}
}
