package middleware

import (
	"context"
	"net/http"
	"shopApi/model"
	"shopApi/token"
	"shopApi/util"
)

func AuthMiddleware(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		jwToken := r.Header.Get("Authorization")
		if jwToken == "" {
			util.Respond(w, &model.Response{Error: "no JSON Web Token provided (set header 'Authorization')"}, http.StatusBadRequest)
			return
		}

		isAuthorized, claims := token.ParseAndVerifyJwt(jwToken)
		if !isAuthorized {
			util.Respond(w, &model.Response{Error: "provided token is not valid"}, http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
