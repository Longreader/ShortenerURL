package config

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
)

func GetAddress() string {

	ParseFlagIfDidnt()

	logrus.Info("Server address configuration")

	ServerAddress := os.Getenv("SERVER_ADDRESS")

	if ServerAddress == "" {
		ServerAddress = *Flags.a
		logrus.Debug("Server address comes from flag and equal", ServerAddress)
	}
	return ServerAddress
}

func GetURL() string {

	ParseFlagIfDidnt()

	logrus.Info("Base URL configuration")

	BaseURL := os.Getenv("BASE_URL")

	if BaseURL == "" {
		BaseURL = *Flags.b
		logrus.Debug("Base URL comes from flag and equal", BaseURL)
	}
	if BaseURL[len(BaseURL)-1:] != "/" {
		BaseURL += "/"
	}
	return BaseURL
}

func GetStoragePath() string {

	ParseFlagIfDidnt()

	logrus.Info("File Storage Path configuration")

	fileName := os.Getenv("FILE_STORAGE_PATH")

	if fileName == "" {
		fileName = *Flags.f
		logrus.Debug("Path comes from flag and equal", fileName)
	}
	return fileName
}

func ParseFlagIfDidnt() {

	if flag.Parsed() {
		logrus.Debug("Plag already parsed")
		return
	}

	flag.Parse()
	logrus.Info("Flag parsed")
}
