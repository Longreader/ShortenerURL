package storage

import (
	"io"
	"sync"

	"github.com/sirupsen/logrus"
)

type Config struct {
	StoragePath string
}

type Storage struct {
	sync.Mutex  // mutex for lock
	storagePath string
	storage     map[string]string
}

func New(cfg Config) (*Storage, error) {
	logrus.Debug("New Storage")
	defer logrus.Debug("New storage created")
	path := cfg.StoragePath
	if path == "" {
		return &Storage{
			storage:     make(map[string]string),
			storagePath: path,
		}, nil
	} else {
		logrus.Debug("New storage creating from file")
		consumer, err := NewConsumer(path)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		defer consumer.Close()
		st := Storage{
			storage:     make(map[string]string),
			storagePath: path,
		}
		for {
			readItem, err := consumer.ReadURL()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					logrus.Error((err))
					return nil, err
				}
			}
			logrus.Debugf("Short url %s long url %s", readItem.ShortURL, readItem.LongURL)
			st.storage[readItem.ShortURL] = readItem.LongURL
		}
		return &st, nil
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

func (st *Storage) Set(key, value string) error {
	st.Lock()
	defer st.Unlock()
	st.set(key, value)
	path := st.storagePath
	if path != "" {
		producer, err := NewProducer(path)
		if err != nil {
			logrus.Error(err)
			return err
		}
		defer producer.Close()
		st := StorageItem{
			ShortURL: key,
			LongURL:  value,
		}
		if err := producer.WriteURL(&st); err != nil {
			logrus.Error(err)
			return err
		}
	}
	return nil

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

// func (st *Storage) getall() {
// 	for r, v := range st.storage {
// 		fmt.Println(r, v)
// 	}
// }

// func (st *Storage) GetAll() {
// 	st.Lock()
// 	defer st.Unlock()
// 	st.getall()
// }
