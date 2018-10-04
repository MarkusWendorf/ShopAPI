package routes

import (
	"github.com/go-chi/chi"
	"math"
	"net/http"
	"shopApi/database"
	"shopApi/model"
	"shopApi/util"
	"shopApi/util/urlQueryParser"
)

func (h *Handlers) GetCategoryNames(w http.ResponseWriter, r *http.Request) {

	categories := h.db.GetCategoryNames()

	util.Respond(w, &model.Response{Data: categories}, http.StatusOK)
}

func (h *Handlers) GetProduct(w http.ResponseWriter, r *http.Request) {

	paramId := chi.URLParam(r, "id")
	if paramId == "" {
		errorMessage := "provided id is empty"
		util.Respond(w, &model.Response{Error: errorMessage}, http.StatusBadRequest)
		return
	}

	product, err := h.db.GetProduct(paramId)
	if err != nil {
		util.Respond(w, &model.Response{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	util.Respond(w, &model.Response{Data: product}, http.StatusOK)
}

func (h *Handlers) QueryProducts(w http.ResponseWriter, r *http.Request) {

	parser := urlQueryParser.New(r.URL.Query())

	name := parser.GetString("name", "")
	category := parser.GetString("category", "")

	// convert from euro to cents
	from := parser.GetInt("priceFrom", 0) * 100
	to := parser.GetInt("priceTo", math.MaxInt32) * 100

	page := parser.GetInt("page", 1)

	query := new(database.QueryBuilder).
		Name(name).
		Category(category).
		Price(from, to).
		Build()

	products, isLast, err := h.db.GetNthPage(query, page)
	if err != nil {
		util.Respond(w, &model.Response{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	links := map[string]interface{}{
		"isLast": isLast,
		"page":   page,
	}

	util.Respond(w, &model.Response{Data: products, Links: links}, http.StatusOK)
}
