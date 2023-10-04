package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress   string
	ServerBaseURL   string
	FileStoragePath string
	DatabaseDSN     string
	CookieKey       []byte
}

// NewConfig - конструктор для Config, сам получит и запишет значения.
//
// Приоритет (меньше - приоритетнее):
//  0. аргументы командной строки
//  1. env-переменные
//  2. JSON-файл с конфигурацией
func NewConfig() Config {
	cfg := Config{
		ServerAddress:   ":8080",
		ServerBaseURL:   "http://localhost:8080",
		FileStoragePath: "",
		CookieKey:       []byte("YandexPracticum"),
	}

	cfg.loadEnv()
	cfg.loadArgs()

	return cfg
}

func (cfg *Config) loadEnv() {
	if s, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		cfg.ServerAddress = s
	}

	if s, ok := os.LookupEnv("BASE_URL"); ok {
		cfg.ServerBaseURL = s
	}

	if s, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		cfg.FileStoragePath = s
	}

	if s, ok := os.LookupEnv("DATABASE_DSN"); ok {
		cfg.DatabaseDSN = s
	}
}

func (cfg *Config) loadArgs() {
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "server address")
	flag.StringVar(&cfg.ServerBaseURL, "b", cfg.ServerBaseURL, "server base url")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "file storage path")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "database data sourse name")

	flag.Parse()
}
