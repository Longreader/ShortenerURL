package config

import (
	"os"
)

func GetAddress() string {

	ServerAddress := os.Getenv("SERVER_ADDRESS")

	if ServerAddress == "" {
		ServerAddress = "127.0.0.1:8080"
	}
	return ServerAddress
}

func GetURL() string {
	BaseURL := os.Getenv("BASE_URL")

	if BaseURL == "" {
		BaseURL = "http://127.0.0.1:8080/"
	}
	if BaseURL[len(BaseURL)-1:] != "/" {
		BaseURL += "/"
	}
	return BaseURL
}
