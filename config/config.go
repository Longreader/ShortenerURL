package config

import (
	"os"
)

func GetAddress() string {

	ServerAddress := os.Getenv("SERVER_ADDRESS")

	if ServerAddress == "" {
		ServerAddress = Flags.a
	}
	return ServerAddress
}

func GetURL() string {
	BaseURL := os.Getenv("BASE_URL")

	if BaseURL == "" {
		BaseURL = Flags.b
	}
	if BaseURL[len(BaseURL)-1:] != "/" {
		BaseURL += "/"
	}
	return BaseURL
}

func GetStoragePath() string {
	fileName := os.Getenv("FILE_STORAGE_PATH")
	if fileName == "" {
		fileName = Flags.f
	}
	return fileName
}
