package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/Longreader/go-shortener-url.git/internal/storage"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestHandler_ShortenerURLHandler(t *testing.T) {
	type fields struct {
		StoragePath string
		BaseURL     string
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	type response struct {
		StatusCode  int
		ContentType string
		BodyPattern string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		response response
	}{
		{
			name: "POSITIVE TEST #1",
			fields: fields{
				StoragePath: "",
				BaseURL:     "http//:localhost:8000/",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "http://localhost:8000/", bytes.NewBuffer([]byte("vk.com"))),
			},
			response: response{
				StatusCode:  201,
				ContentType: "",
				BodyPattern: `\w\d\w\d\w\d\w`,
			},
		},
		{
			name: "NEGATIVE TEST #1",
			fields: fields{
				StoragePath: "",
				BaseURL:     "http//:localhost:8000/",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "http://localhost:8000/", bytes.NewBuffer([]byte("vk.com"))),
			},
			response: response{
				StatusCode:  400,
				ContentType: "text/plain; charset=utf-8",
				BodyPattern: `Bad request`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			Store, err := storage.New(
				storage.Config{
					StoragePath: "",
				},
			)

			assert.Nil(t, err)

			h := &Handler{
				Store:   Store,
				BaseURL: tt.fields.BaseURL,
			}

			handler := http.HandlerFunc(h.ShortenerURLHandler)

			handler.ServeHTTP(tt.args.w, tt.args.r)

			resp := tt.args.w.Result()
			defer resp.Body.Close()

			resBody, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)

			fmt.Println(resp.StatusCode)
			fmt.Println(resp.Header.Get("Content-Type"))
			fmt.Println(string(resBody))

			assert.Equal(t, tt.response.StatusCode, resp.StatusCode)
			assert.Equal(t, tt.response.ContentType, resp.Header.Get("Content-Type"))
			assert.Regexp(t, regexp.MustCompile(tt.response.BodyPattern), string(resBody))
		})
	}
}

func TestHandler_IDGetHandler(t *testing.T) {
	type fields struct {
		StoragePath string
		BaseURL     string
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	type response struct {
		StatusCode  int
		ContentType string
		BodyPattern string
	}
	type request struct {
		ShortURL string
		ShortID  string
		LongURL  string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		response response
		request  request
	}{
		{
			name: "POSITIVE TEST #1",
			fields: fields{
				StoragePath: "",
				BaseURL:     "http//:localhost:8000/",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "http://localhost:8000/"+"f3r7h3s", nil),
			},
			response: response{
				StatusCode:  307,
				ContentType: "text/html; charset=UTF-8",
			},
			request: request{
				ShortURL: "http://localhost:8000/f3r7h3s",
				ShortID:  "f3r7h3s",
				LongURL:  "https://tproger.ru/articles/puteshestvie-v-golang-regexp",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			Store, err := storage.New(
				storage.Config{
					StoragePath: tt.fields.StoragePath,
				},
			)

			assert.Nil(t, err)

			h := &Handler{
				Store:   Store,
				BaseURL: tt.fields.BaseURL,
			}

			h.Store.Set(tt.request.ShortID, tt.request.LongURL)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.request.ShortID)

			tt.args.r = tt.args.r.WithContext(context.WithValue(tt.args.r.Context(), chi.RouteCtxKey, rctx))

			handler := http.HandlerFunc(h.IDGetHandler)

			handler.ServeHTTP(tt.args.w, tt.args.r)

			resp := tt.args.w.Result()

			defer resp.Body.Close()

			resBody, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)

			fmt.Println(resp.StatusCode)
			fmt.Println(resp.Header.Get("Content-Type"))
			fmt.Println(string(resBody))

			assert.Equal(t, tt.response.StatusCode, resp.StatusCode)
			assert.Equal(t, tt.response.ContentType, resp.Header.Get("Content-Type"))
		})
	}
}

func TestHandler_APIShortenerURLHandler(t *testing.T) {
	type fields struct {
		StoragePath string
		BaseURL     string
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	type response struct {
		StatusCode  int
		ContentType string
		BodyPattern string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		response response
	}{
		{
			name: "POSITIVE TEST #1",
			fields: fields{
				StoragePath: "",
				BaseURL:     "http//:localhost:8000/",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "http://localhost:8000/api/shorten", bytes.NewBuffer([]byte(`{"url":"https://github.com/Longreader/ShortenerURL/"}`))),
			},
			response: response{
				StatusCode:  201,
				ContentType: "application/json",
				BodyPattern: `\w\d\w\d\w\d\w`,
			},
		},
		{
			name: "POSITIVE TEST #2",
			fields: fields{
				StoragePath: "",
				BaseURL:     "http//:localhost:8000/",
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "http://localhost:8000/api/shorten", bytes.NewBuffer([]byte(`{"flip":"https://github.com/Longreader/ShortenerURL/"}`))),
			},
			response: response{
				StatusCode:  201,
				ContentType: "application/json",
				BodyPattern: `\w\d\w\d\w\d\w`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			Store, err := storage.New(
				storage.Config{
					StoragePath: tt.fields.StoragePath,
				},
			)

			assert.Nil(t, err)

			h := &Handler{
				Store:   Store,
				BaseURL: tt.fields.BaseURL,
			}

			tt.args.r.Header.Set("Content-Type", "application/json")

			handler := http.HandlerFunc(h.APIShortenerURLHandler)

			handler.ServeHTTP(tt.args.w, tt.args.r)

			resp := tt.args.w.Result()
			defer resp.Body.Close()

			resBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			body := make(map[string]string)
			err = json.Unmarshal(resBody, &body)

			fmt.Printf(body["url"])

			assert.Nil(t, err)

			assert.Equal(t, tt.response.StatusCode, resp.StatusCode)
			assert.Equal(t, tt.response.ContentType, resp.Header.Get("Content-Type"))
			assert.Regexp(t, regexp.MustCompile(tt.fields.BaseURL+tt.response.BodyPattern), body["result"])
		})
	}
}
