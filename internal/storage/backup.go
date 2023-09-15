package storage

import (
	"encoding/json"
	"os"
)

type StorageItem struct {
	LongURL  string `json:"long_URL"`
	ShortURL string `json:"short_URL"`
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

type consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewProducer(fileName string) (*producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) WriteURL(st *StorageItem) error {
	return p.encoder.Encode(&st)
}

func (p *producer) Close() error {
	return p.file.Close()
}

func NewConsumer(fileName string) (*consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *consumer) ReadURL() (*StorageItem, error) {
	st := &StorageItem{}
	if err := c.decoder.Decode(&st); err != nil {
		return nil, err
	}
	return st, nil
}

func (c *consumer) Close() error {
	return c.file.Close()
}
