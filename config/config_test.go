package config

import (
	"testing"
)

func TestGetValue(t *testing.T) {
	type want struct {
		ServerAddress   string
		ServerBaseURL   string
		FileStoragePath string
	}

	tests := []struct {
		name string
		foo  func() Config
		want want
	}{
		{
			name: "positive test #1 loadEnv",
			foo:  NewConfig,
			want: want{
				ServerAddress:   ":8080",
				ServerBaseURL:   "http://localhost:8080",
				FileStoragePath: "",
			},
		},
		// {
		// 	name: "positive test #2 loadArgs",
		// 	foo:  NewConfig,
		// 	want: want{
		// 		ServerAddress:   ":8080",
		// 		ServerBaseURL:   "http://localhost:8080",
		// 		FileStoragePath: "",
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.foo()
			if result.FileStoragePath != tt.want.FileStoragePath {
				t.Errorf("Expected value code %s, got %s", tt.want.FileStoragePath, result.FileStoragePath)
			}
			if result.ServerAddress != tt.want.ServerAddress {
				t.Errorf("Expected value code %s, got %s", tt.want.ServerAddress, result.ServerAddress)
			}
			if result.ServerBaseURL != tt.want.ServerBaseURL {
				t.Errorf("Expected value code %s, got %s", tt.want.ServerBaseURL, result.ServerBaseURL)
			}
		})
	}

}
