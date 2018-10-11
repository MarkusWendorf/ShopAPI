package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"shopApi/model"
	"shopApi/util"
)

type CartRequestBody struct {
	Items  []model.CartItem `json:"items"`
	CartId string           `json:"cart_id"`
}

func (h Handlers) Cart(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var cartRequest CartRequestBody
	err := json.NewDecoder(r.Body).Decode(&cartRequest)
	if err != nil {
		util.Respond(w, &model.Response{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	items := cartRequest.Items
	// map[productId]quantity
	quantities := map[string]int{}

	for _, item := range items {
		id := item.Product.Id
		quantities[id] = item.Quantity
	}

	claims, ok := r.Context().Value("claims").(model.JwtClaims)
	if !ok {
		log.Fatal("couldn't get claims from request context")
	}

	validatedCart, err := h.db.ValidateCart(claims.UserId, cartRequest.CartId, quantities)
	if err != nil {
		util.Respond(w, &model.Response{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	util.Respond(w, &model.Response{Data: validatedCart}, http.StatusOK)
}
