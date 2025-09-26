package config

import (
	"backend/pkg/logger"
	"context"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

type Config struct {
	Port    int    `env:"PORT"`
	Env     string `env:"ENV"`
	SQLPath string `env:"SQLPATH"`
}

func loadEnvVariablesFromFile(envPath ...string) error {
	var err error
	if len(envPath) == 0 {
		err = godotenv.Overload()
	} else {
		err = godotenv.Overload(envPath...)
	}

	return err
}

func MustGetConfig(ctx context.Context, log logger.Logger) *Config {
	if err := loadEnvVariablesFromFile(); err != nil {
		log.Warn("failed to load env variables", zap.Error(err))
	}

	config := &Config{}
	if err := envconfig.Process(ctx, config); err != nil {
		log.Warn("failed to load config", zap.Error(err))
		return &Config{
			Port:    8080,
			Env:     "dev",
			SQLPath: "./storage/storage.db", // TODO: Сделать postgres
		}
	}

	log.Info("initialized config", *config)
	return config
}
