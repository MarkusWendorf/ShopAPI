package middleware

import (
	"net/http"
)

type CorsConfig struct {
	AllowedOrigins   []string
	AllowedHeaders   string
	AllowCredentials bool
}

func CorsMiddleware(config CorsConfig) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		fn := func(w http.ResponseWriter, r *http.Request) {

			origin := r.Header.Get("origin")

			if !isValidOrigin(config.AllowedOrigins, origin) {
				// skip; dont set cors headers
				next.ServeHTTP(w, r)
				return
			}

			headers := w.Header()
			headers.Add("Access-Control-Allow-Origin", origin)

			if config.AllowCredentials {
				headers.Add("Access-Control-Allow-Credentials", "true")
			}

			if config.AllowedHeaders != "" {
				headers.Add("Access-Control-Allow-Headers", config.AllowedHeaders)
			}

			if r.Method == http.MethodOptions {
				// no further handling for method OPTIONS
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func isValidOrigin(allowedOrigins []string, origin string) bool {

	for _, allowed := range allowedOrigins {
		if origin == allowed || allowed == "*" {
			return true
		}
	}

	return false
}
