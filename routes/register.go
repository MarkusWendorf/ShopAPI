package routes

import (
	"shopApi/model"
	"shopApi/util"
	"encoding/json"
	"net/http"
)

type RegisterRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h Handlers) Register(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var loginRequest RegisterRequestBody
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		util.Respond(w, &model.Response{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	_, err = h.db.NewUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		util.Respond(w, &model.Response{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	util.Respond(w, &model.Response{}, http.StatusOK)
}
