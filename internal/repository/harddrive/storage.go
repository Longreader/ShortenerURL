package harddrive

import (
	"bufio"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/Longreader/go-shortener-url.git/internal/repository"
	"github.com/Longreader/go-shortener-url.git/internal/repository/memory"
	"github.com/google/uuid"
)

type FileStorage struct {
	file *os.File
	memory.MemoryStorage
	fileMutex sync.Mutex
}

func NewFileStorage(file *os.File) (*FileStorage, error) {

	st := &FileStorage{
		file: file,
	}

	st.LinksStorage = make(map[repository.URL]repository.ID)
	st.UserLinkStorage = make(map[repository.ID]repository.LinkData)

	err := st.load()
	if err != nil {
		return nil, err
	}

	return st, nil
}

func (st *FileStorage) Set(ctx context.Context, url repository.URL, user repository.User) (id repository.ID, err error) {
	id, err = st.SetLink(url, user)
	if err != nil {
		return
	}
	err = st.write(fmt.Sprintf("%s,%s,%s", id, user.String(), base64.StdEncoding.EncodeToString([]byte(url))))
	return
}

func (st *FileStorage) load() error {

	st.Lock()
	defer st.Unlock()

	reader := bufio.NewReader(st.file)

	i := 0

	for {
		bytes, err := reader.ReadBytes('\n')
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		line := strings.Trim(string(bytes), "\n")
		splitted := strings.Split(line, ",")

		id := splitted[0]
		user, err := uuid.Parse(splitted[1])
		if err != nil {
			return repository.ErrUnableParseUser
		}

		var data []byte
		data, err = base64.StdEncoding.DecodeString(splitted[3])
		if err != nil {
			return repository.ErrUnableDecodeURL
		}
		url := repository.URL(data)

		st.UserLinkStorage[id] = repository.LinkData{
			URL:  url,
			User: user,
		}
		st.LinksStorage[url] = id

		i++
	}
	return nil
}

func (st *FileStorage) write(data string) error {
	st.fileMutex.Lock()
	defer st.fileMutex.Unlock()

	_, err := st.file.Write([]byte(data + "\n"))
	if err != nil {
		return err
	}
	return nil
}
