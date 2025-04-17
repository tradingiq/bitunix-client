package samples

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Cfg struct {
	ApiKey    string `env:"API_KEY"`
	SecretKey string `env:"SECRET_KEY"`
}

var Config Cfg

func init() {
	log.SetLevel(log.DebugLevel)
	err := godotenv.Load()
	if err != nil {
		log.Warn("Error loading .env file")
	}

	err = env.Parse(&Config)
	Config, err = env.ParseAs[Cfg]()

	if err != nil {
		log.WithError(err).Fatal("parse env vars")
	}

}
