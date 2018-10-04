package routes

import (
	"shopApi/model"
	"shopApi/util"
	"encoding/json"
	"net/http"
)

type OrderRequestBody struct {
	CartId  string                   `json:"cart_id"`
	Address model.AddressInformation `json:"address"`
}

func (h Handlers) Order(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var cartRequest OrderRequestBody
	err := json.NewDecoder(r.Body).Decode(&cartRequest)
	if err != nil {
		util.Respond(w, &model.Response{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	cart, err := h.db.GetCart(cartRequest.CartId)
	if err != nil {
		util.Respond(w, &model.Response{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	order := model.Order{Cart: cart, Address: cartRequest.Address}

	err = h.db.PutOrder(&order)
	if err != nil {
		util.Respond(w, &model.Response{Error: err.Error()}, http.StatusInternalServerError)
		return
	}

	util.Respond(w, &model.Response{Data: order.Id}, http.StatusOK)
}
