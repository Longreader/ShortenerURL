package config

import (
	"log"
	"os"
)

var SERVER_ADDRESS = os.Getenv("SERVER_ADDRESS")

var BASE_URL = os.Getenv("BASE_URL")

func Setup() {
	if SERVER_ADDRESS == "" {
		SERVER_ADDRESS = "127.0.0.1:8080"
	}

	if BASE_URL == "" {
		BASE_URL = "http://127.0.0.1:8080/"
	}
	log.Println("Server setup configuration completed")
}
