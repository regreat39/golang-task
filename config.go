package main

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port           string
	DBPath         string
	NextTakingsMin int
}

var defaultCfg = Config{
	Port:           "8080",
	DBPath:         "medicine.db",
	NextTakingsMin: 60,
}

func NewFromEnv() (Config, error) {
	cfg := defaultCfg

	if val, found := os.LookupEnv("PORT"); found && val != "" {
		cfg.Port = val
	}

	if val, found := os.LookupEnv("DB_PATH"); found && val != "" {
		cfg.DBPath = val
	}

	if val, found := os.LookupEnv("NEXT_TAKINGS_MIN"); found && val != "" {
		num, err := strconv.Atoi(val)
		if err != nil {
			return Config{}, fmt.Errorf("NEXT_TAKINGS_MIN должно быть числом, получено %q", val)
		}
		cfg.NextTakingsMin = num
	}

	if cfg.Port == "" {
		return Config{}, fmt.Errorf("PORT не должен быть пустым")
	}

	if cfg.DBPath == "" {
		return Config{}, fmt.Errorf("DB_PATH не должен быть пустым")
	}

	if cfg.NextTakingsMin <= 0 {
		return Config{}, fmt.Errorf("NEXT_TAKINGS_MIN должен быть положительным, текущее значение: %d", cfg.NextTakingsMin)
	}

	return cfg, nil
}
