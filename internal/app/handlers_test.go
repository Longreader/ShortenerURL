package app_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/Longreader/go-shortener-url.git/internal/app"
	"github.com/gorilla/mux"
)

func TestPostEndpoint(t *testing.T) {

	// var baseURL string = "http://localhost:8080/"

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

			// if string(resBody) != tt.want.response {
			// 	t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			// }
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

	var baseURL = "http://127.0.0.1:8080/"

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
				code:     307,
				response: "",
			},
		},
		{
			name:      "positive test #2 GET",
			value:     "https://practicum.yandex.ru/",
			key:       "m8J7h9R",
			searchKey: "f0f0f0f",
			want: want{
				code:     400,
				response: "Bad request",
			},
		},
	}

	for _, tt := range tests {
		// Запускаем хендлеры
		t.Run(tt.name, func(t *testing.T) {

			app.Store.Set(tt.key, tt.value)

			request := httptest.NewRequest(http.MethodGet, baseURL+tt.searchKey, nil)

			// value, okey := app.Store.Get(tt.searchKey)
			// if okey != true {
			// 	fmt.Printf("%s", "Error occured")
			// }
			// fmt.Printf("%s", value)

			val := map[string]string{
				"id": tt.searchKey,
			}

			// выставляем параметр id в url vars
			request = mux.SetURLVars(request, val)

			// создаём новый Recorder
			w := httptest.NewRecorder()

			// определяем хендлер
			h := http.HandlerFunc(app.IDGetHandler)

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

			if strings.Trim(string(resBody), "\n \t") != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}

		})
	}
}
