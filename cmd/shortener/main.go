package main

import (
	"flag"
	"net/http"

	"github.com/Longreader/go-shortener-url.git/internal/app"
	"github.com/Longreader/go-shortener-url.git/internal/config"
	"github.com/Longreader/go-shortener-url.git/internal/storage"
	"github.com/sirupsen/logrus"
)

func main() {

	logrus.StandardLogger().Level = logrus.DebugLevel
	logrus.SetFormatter(&logrus.JSONFormatter{})

	flag.Parse()

	cfg := config.NewConfig()

	db, err := storage.New(storage.Config{
		StoragePath: cfg.StoragePath,
	})

	if err != nil {
		logrus.Fatal("Error database connection: ", err)
	}

	h := app.NewHandler(
		db,
		cfg.BaseURL,
	)

	r := h.InitRouter()

	http.Handle("/", r)

	logrus.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
