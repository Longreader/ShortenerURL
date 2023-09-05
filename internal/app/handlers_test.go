package app_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"regexp"
	"testing"

	"github.com/Longreader/go-shortener-url.git/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", app.IDGetHandler)
		r.Post("/", app.ShortenerURLHandler)
	})
	return r
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}

func TestPostEndpoint(t *testing.T) {

	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name  string
		value string
		want  want
	}{
		{
			name:  "positive test #1 POST",
			value: "https://practicum.yandex.ru/",
			want: want{
				code:     201,
				response: `http://127.0.0.1:8080/`,
			},
		},
		{
			name:  "positive test #2 POST",
			value: "https://tproger.ru/articles/puteshestvie-v-golang-regexp/",
			want: want{
				code:     201,
				response: `http://127.0.0.1:8080/`,
			},
		},
	}

	for _, tt := range tests {
		// Запускаем хендлеры
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.value, nil)

			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(app.ShortenerURLHandler)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			// получаем и проверяем тело запроса
			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			matched, _ := regexp.MatchString((tt.want.response + `\w+`), string(resBody))
			if !matched {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())

			}

			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}

		})
	}
}

func TestGetEndpoint(t *testing.T) {

	type want struct {
		code     int
		response string
	}

	tests := []struct {
		name      string
		key       string
		searchKey string
		value     string
		want      want
	}{
		{
			name:      "positive test #1 GET",
			value:     "https://practicum.yandex.ru/",
			key:       "x4G3v6K",
			searchKey: "x4G3v6K",
			want: want{
				code: 200,
			},
		},
		{
			name:      "positive test #2 GET",
			value:     "https://practicum.yandex.ru/",
			key:       "m8J7h9R",
			searchKey: "f0f0f0f",
			want: want{
				code:     400,
				response: "Bad request\n",
			},
		},
	}

	for _, tt := range tests {

		r := NewRouter()
		ts := httptest.NewServer(r)
		defer ts.Close()

		app.Store.Set(tt.key, tt.value)

		statusCode, body := testRequest(t, ts, "GET", "/"+tt.searchKey)
		assert.Equal(t, tt.want.code, statusCode)
		if statusCode != http.StatusOK {
			assert.Equal(t, tt.want.response, body)
		}
	}
}

func TestAPIPostEndpoint(t *testing.T) {

	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name  string
		value string
		want  want
	}{
		// {
		// 	name:  "positive test #1 POST",
		// 	value: "https://practicum.yandex.ru/",
		// 	want: want{
		// 		code:        201,
		// 		response:    ,
		// 		contentType: "application/json",
		// 	},
		// },
		{
			name:  "negative test #1 POST",
			value: "https://tproger.ru/articles/puteshestvie-v-golang-regexp/",
			want: want{
				code:        406,
				response:    "Bad agent request",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		// Запускаем хендлеры
		t.Run(tt.name, func(t *testing.T) {

			// создаем словарь данных на отправку
			var data = make(map[string]string)
			data["url"] = tt.value

			// создание буфера для записи json
			buf := bytes.NewBuffer([]byte{})
			// кодирование данных в буфер
			encoder := json.NewEncoder(buf)
			encoder.Encode(data)

			request := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8080/api/shorten", bytes.NewBuffer(buf.Bytes()))

			request.Header.Add("Content-Type", tt.want.contentType)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(app.APIShortenerURLHandler)
			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			// получаем и проверяем тело запроса
			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			// if string(resBody) != tt.want.response {
			// 	t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			// }

			resBodyStr := strings.Trim(string(resBody), "\n")

			matched, _ := regexp.MatchString(tt.want.response, resBodyStr)
			if !matched {
				t.Errorf("Expected body %s, got %s and %s word?", tt.want.response, w.Body.String(), string(resBody))

			}

			// заголовок ответа
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}

		})
	}
}
