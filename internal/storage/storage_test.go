package storage_test

import (
	"os"
	"testing"

	"github.com/Longreader/go-shortener-url.git/internal/storage"
)

func TestStorage(t *testing.T) {
	var valueURL = "value"
	var keyURL = "key"
	os.Setenv("SERVER_ADDRESS", "127.0.0.1:8080")
	os.Setenv("BASE_URL", "http://127.0.0.1:8080/")
	os.Setenv("FILE_STORAGE_PATH", "log.log")
	store := storage.New()
	store.Set(keyURL, valueURL)
	if storeVal, _ := store.Get(keyURL); storeVal != valueURL {
		t.Errorf("Expected value %s, got %s", storeVal, valueURL)
	}
}
