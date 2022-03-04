package config

import (
	"gade/srv-gade-point/logger"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// check .env file existence
	envPath := ".env"

	if os.Getenv(`ENV_PATH`) != "" {
		envPath = os.Getenv(`ENV_PATH`)
	}

	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return
	}

	err := godotenv.Load(envPath)

	if err != nil {
		logger.Make(nil, nil).Fatal(err)
	}
}
