package config

import (
	"flag"

	"github.com/sirupsen/logrus"
)

var Flags struct {
	a *string
	b *string
	f *string
}

func init() {
	Flags.a = flag.String("a", "127.0.0.1:8080", "SERVER_ADDRESS")

	Flags.b = flag.String("b", "http://127.0.0.1:8080", "BASE_URL")

	Flags.f = flag.String("f", "", "FILE_STORAGE_PATH")

	logrus.Info("Init Flags")
}
