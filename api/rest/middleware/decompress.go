package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// Decompress middleware
func Decompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var b bytes.Buffer
		if _, err = b.ReadFrom(gz); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		nr := io.NopCloser(bytes.NewReader(b.Bytes()))
		rb, _ := http.NewRequest(r.Method, r.RequestURI, nr)
		_ = r.Body.Close()
		next.ServeHTTP(w, rb)
	})
}
