package config

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

func TestGetValue(t *testing.T) {
	type want struct {
		value string
	}

	tests := []struct {
		name string
		foo  func() string
		want want
	}{
		{
			name: "positive test #1 GetAddress",
			foo:  getAddress,
			want: want{
				value: "127.0.0.1:8080",
			},
		},
		{
			name: "positive test #2 GetURL",
			foo:  getURL,
			want: want{
				value: "http://127.0.0.1:8080/",
			},
		},
		{
			name: "positive test #2 GetPath",
			foo:  getStoragePath,
			want: want{
				value: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.foo()
			if result != tt.want.value {
				t.Errorf("Expected value code %s, got %s", tt.want.value, result)
			}
		})
	}

}
