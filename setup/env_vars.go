package setup

import (
	"github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
)

func LoadEnvironmentalVariables(logger hclog.Logger) {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Warn("Could not detect .env file")
	}
}
