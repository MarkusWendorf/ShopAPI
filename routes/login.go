package routes

import (
	"encoding/json"
	"net/http"
	"shopApi/model"
	"shopApi/token"
	"shopApi/util"
)

type LoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponseBody struct {
	Token string `json:"token"`
}

func (h Handlers) Login(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var loginRequest LoginRequestBody
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		util.Respond(w, &model.Response{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	user, err := h.db.GetUser(loginRequest.Email)
	if err != nil {
		util.Respond(w, &model.Response{Error: err.Error()}, http.StatusForbidden)
		return
	}

	jwToken, err := token.CreateJWT(user)
	loginResponse := LoginResponseBody{Token: jwToken}

	util.Respond(w, &model.Response{Data: loginResponse}, http.StatusOK)
}
