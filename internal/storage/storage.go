package storage

import (
	"sync"
)

type Storage struct {
	sync.Mutex // mutex for lock
	storage    map[string]string
}

func New() *Storage {
	return &Storage{
		storage: make(map[string]string),
	}
}

func (st *Storage) set(key string, value string) {
	st.storage[key] = value
}

func (st *Storage) get(key string) (string, bool) {
	// Почему ок возвращает ошибку
	// но при этом возвращается верный результат
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
