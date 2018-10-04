package util

import (
	"encoding/json"
	"log"
	"net/http"
	"shopApi/model"
)

func Respond(w http.ResponseWriter, data *model.Response, status int) {

	w.Header().Add("Content-Type", "application/json")

	b, err := json.Marshal(data)
	if err != nil {

		errResponse := model.Response{Error: err.Error()}
		b, err = json.Marshal(errResponse)
		if err != nil {
			log.Fatalf("couldn't marshal error response for error: %s", errResponse.Error)
		}

		status = http.StatusInternalServerError
	}

	w.WriteHeader(status)
	w.Write(b)
}
