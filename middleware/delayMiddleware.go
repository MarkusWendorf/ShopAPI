package middleware

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

func DelayMiddleware(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		delay := rand.Intn(3)
		log.Println("Delay in ms:", delay)

		time.Sleep(time.Duration(delay) * time.Millisecond)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}