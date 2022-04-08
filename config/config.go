package config

import (
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func LoadEnv() {
	envPath := ".env"

	if os.Getenv(`APP_PATH`) != "" {
		envPath = os.Getenv(`APP_PATH`) + "/" + envPath
	}

	// check .env file existence
	_, err := os.Stat(envPath)

	if os.IsNotExist(err) {
		logger.Make(nil).Debug(models.ErrEnvFileNF)
		return
	}

	_ = godotenv.Load(envPath)
}

func LoadTestData() {
	mockDataPath := os.Getenv(`APP_PATH`) + "/mocks/data"
	items, err := ioutil.ReadDir(mockDataPath)

	if err != nil {
		logger.Make(nil).Fatal(err)
	}

	viper.AddConfigPath(mockDataPath)
	_ = viper.ReadInConfig()

	mockDataFile := mockDataPath
	for _, item := range items {
		if item.IsDir() {
			// currently we does not expect any dir inside
			continue
		}

		// load all existed yaml file
		mockDataFile = mockDataFile + "/" + item.Name()
		viper.SetConfigName(strings.TrimSuffix(filepath.Base(item.Name()), filepath.Ext(item.Name())))
		viper.AddConfigPath(mockDataFile)
		err = viper.MergeInConfig()

		if err != nil {
			logger.Make(nil).Fatal(err)
		}
	}
}
