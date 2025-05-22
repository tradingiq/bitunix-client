package samples

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Cfg struct {
	ApiKey    string `env:"API_KEY"`
	SecretKey string `env:"SECRET_KEY"`
}

var Config Cfg

func init() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	err := godotenv.Load()
	if err != nil {
		logger.Warn("Error loading .env file")
	}

	err = env.Parse(&Config)
	Config, err = env.ParseAs[Cfg]()

	if err != nil {
		logger.Fatal("parse env vars", zap.Error(err))
	}
}
