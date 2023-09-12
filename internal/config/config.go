package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Config struct {
	BaseURL       string
	ServerAddress string
	StoragePath   string
}

func NewConfig() *Config {
	BaseURL := getURL()
	ServerAddress := getAddress()
	StoragePath := getStoragePath()
	return &Config{
		BaseURL:       BaseURL,
		ServerAddress: ServerAddress,
		StoragePath:   StoragePath,
	}
}

func getAddress() string {

	logrus.Info("Server address configuration")

	ServerAddress := os.Getenv("SERVER_ADDRESS")

	if ServerAddress == "" {
		ServerAddress = *Flags.a
		logrus.Debug("Server address comes from flag and equal ", ServerAddress)
	} else {
		logrus.Debug("Server address comes from env and equal ", ServerAddress)
	}
	return ServerAddress
}

func getURL() string {

	logrus.Info("Base URL configuration")

	BaseURL := os.Getenv("BASE_URL")

	if BaseURL == "" {
		BaseURL = *Flags.b
		logrus.Debug("Base URL comes from flag and equal ", BaseURL)
	} else {
		logrus.Debug("Base URL comes from env and equal ", BaseURL)
	}
	if BaseURL[len(BaseURL)-1:] != "/" {
		BaseURL += "/"
	}
	return BaseURL
}

func getStoragePath() string {

	logrus.Info("File Storage Path configuration")

	Path := os.Getenv("FILE_STORAGE_PATH")

	if Path == "" {
		Path = *Flags.f
		logrus.Debug("Path comes from flag and equal ", Path)
	} else {
		logrus.Debug("Path comes from env and equal ", Path)
	}
	return Path
}
