package storage

import (
	"io"
	"log"
	"sync"

	"github.com/Longreader/go-shortener-url.git/config"
)

type Storage struct {
	sync.Mutex // mutex for lock
	storage    map[string]string
}

func New() *Storage {
	fileName := config.GetStoragePath()
	// log.Printf("The fileName is %s", fileName)
	if fileName == "" {
		return &Storage{
			storage: make(map[string]string),
		}
	} else {
		consumer, err := NewConsumer(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer consumer.Close()
		st := Storage{
			storage: make(map[string]string),
		}
		for {
			readItem, err := consumer.ReadURL()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					log.Fatal(err)
				}
			}
			st.storage[readItem.ShortURL] = readItem.LongURL
		}
		return &st
	}
}

func (st *Storage) set(key string, value string) {
	st.storage[key] = value
}

func (st *Storage) get(key string) (string, bool) {
	if st.count() > 0 {
		item, ok := st.storage[key]
		if !ok {
			return "/", false
		}
		return item, ok
	}
	return "/", false
}

func (st *Storage) count() int {
	return len(st.storage)
}

func (st *Storage) Set(key, value string) {
	st.Lock()
	defer st.Unlock()
	st.set(key, value)
	fileName := config.GetStoragePath()
	log.Printf("The fileName is %s", fileName)
	if fileName != "" {
		produser, err := NewProduser(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer produser.Close()
		st := StorageItem{
			ShortURL: key,
			LongURL:  value,
		}
		if err := produser.WriteURL(&st); err != nil {
			log.Fatal(err)
		}
	}

}

func (st *Storage) Get(key string) (string, bool) {
	st.Lock()
	defer st.Unlock()
	return st.get(key)
}

func (st *Storage) Count() int {
	st.Lock()
	defer st.Unlock()
	return st.count()
}
