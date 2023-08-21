package config

import (
	"os"
)

func GetAddress() string {

	Server_Address := os.Getenv("SERVER_ADDRESS")

	if Server_Address == "" {
		Server_Address = "127.0.0.1:8080"
	}
	return Server_Address
}

func GetURL() string {
	Base_URL := os.Getenv("BASE_URL")

	if Base_URL == "" {
		Base_URL = "http://127.0.0.1:8080/"
	}
	return Base_URL
}
