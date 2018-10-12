package routes

import (
	"github.com/go-chi/chi"
	"net/http"
	"net/url"
	"shopApi/database"
	"shopApi/model"
	"shopApi/util"
	"strconv"
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

func (h *Handlers) Autocomplete(w http.ResponseWriter, r *http.Request) {

	search := chi.URLParam(r, "search")

	matches := h.db.Autocomplete(search)
	util.Respond(w, &model.Response{Data: matches}, http.StatusOK)
}

func (h *Handlers) QueryProducts(w http.ResponseWriter, r *http.Request) {

	values := r.URL.Query()
	query := new(database.QueryBuilder)

	if name := values.Get("name"); name != "" {
		query = query.Name(name)
	}

	if category := values.Get("category"); category != "" {
		query = query.Category(category)
	}

	page, err := strconv.Atoi(values.Get("page"))
	if err != nil {
		page = 1
	}

	from := getInt(values, "priceFrom", 100)
	to := getInt(values, "priceTo", 100)
	query = query.Price(from, to)

	products, last, err := h.db.ExecuteQuery(query, 24, page)
	if err != nil {
		util.Respond(w, &model.Response{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	links := map[string]interface{}{
		"last": last,
		"page": page,
	}

	util.Respond(w, &model.Response{Data: products, Links: links}, http.StatusOK)
}

func getInt(values url.Values, key string, multiplier int) *int {

	val := values.Get(key)
	if val == "" {
		return nil
	}

	toInt, err := strconv.Atoi(val)
	if err != nil {
		return nil
	}

	toInt *= multiplier
	return &toInt
}
