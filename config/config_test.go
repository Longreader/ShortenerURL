package config_test

import (
	"flag"
	"testing"
)

func TestFlags(t *testing.T) {
	val := flag.String("s", ":4040", "PORT")
	flag.Parse()
	if *val != ":4040" {
		t.Errorf("Expected port :4040, got %s", *val)
	}
}
