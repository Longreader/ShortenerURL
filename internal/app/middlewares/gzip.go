package middlewares

import (
	"compress/gzip"
	"net/http"
)

func (m *Middlewares) DecompresMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get(`Content-Encoding`) == `gzip` {

			var reader *gzip.Reader
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				next.ServeHTTP(w, r)
				return
			}
			reader = gz
			defer gz.Close()
			r.Body = reader
		}

		next.ServeHTTP(w, r)

	})
}
