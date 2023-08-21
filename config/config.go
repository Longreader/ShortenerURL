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
	BaseUrl := os.Getenv("BASE_URL")

	if BaseUrl == "" {
		BaseUrl = "http://127.0.0.1:8080/"
	}
	if BaseUrl[len(BaseUrl)-1:] != "/" {
		BaseUrl += "/"
	}
	return BaseUrl
}
