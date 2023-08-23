package config

import "flag"

var Flags struct {
	a string
	b string
	f string
}

func init() {
	sa := flag.String("a", "127.0.0.1:8080", "SERVER_ADDRESS")

	bu := flag.String("b", "http://127.0.0.1:8080/", "BASE_URL")

	fsp := flag.String("f", "", "FILE_STORAGE_PATH")

	flag.Parse()

	Flags.a = *sa
	Flags.b = *bu
	Flags.f = *fsp
}
