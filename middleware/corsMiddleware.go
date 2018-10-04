package middleware

import (
	"net/http"
)

const origin = "http://localhost:8080"

func CorsMiddleware(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		headers := w.Header()
		headers.Add("Access-Control-Allow-Origin", origin)
		headers.Add("Access-Control-Allow-Credentials", "true")
		headers.Add("Access-Control-Allow-Headers", "authorization")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
