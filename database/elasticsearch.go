package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"log"
	"shopApi/model"
	"strings"
)

type QueryBuilder struct {
	queries []elastic.Query
	filters []elastic.Query
}

func (builder *QueryBuilder) Name(pname string) *QueryBuilder {
	matchName := elastic.NewMatchQuery("pname", pname)
	builder.queries = append(builder.queries, matchName)

	return builder
}

func (builder *QueryBuilder) Price(from *int, to *int) *QueryBuilder {
	filterPrice := elastic.NewRangeQuery("price").From(from).To(to)
	builder.filters = append(builder.filters, filterPrice)

	return builder
}

func (builder *QueryBuilder) Category(category string) *QueryBuilder {
	exactCategory := elastic.NewTermQuery("category", strings.ToLower(category))
	builder.queries = append(builder.queries, exactCategory)

	return builder
}

func (builder *QueryBuilder) Execute(client *elastic.Client, pageSize int, page int) ([]model.Product, int, error) {

	if page < 1 {
		return nil, 1, fmt.Errorf("page number should be > 1, got: %d", page)
	}

	query := elastic.NewBoolQuery().
		Must(builder.queries...).
		Filter(builder.filters...)

	search := client.Search("products").Type("product").
		Query(query).
		From((page - 1) * pageSize).
		Size(pageSize)

	return getResults(search, pageSize)
}

func (db *Database) ExecuteQuery(builder *QueryBuilder, pageSize int, page int) ([]model.Product, int, error) {
	return builder.Execute(db.elasticClient, pageSize, page)
}

func (db *Database) Autocomplete(input string) []model.AutocompleteProduct {

	highlight := elastic.NewHighlight().Field("pname").
		PreTags("<b>").
		PostTags("</b>")

	query := elastic.NewMatchQuery("pname", input)
	res, err := db.elasticClient.Search("products").Type("product").
		Query(query).
		Highlight(highlight).
		MinScore(3).
		Do(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	products, err := collectAutocomplete(res.Hits)
	if err != nil {
		log.Fatal(err)
	}

	return products
}



// ===== Helpers =====

func getResults(search *elastic.SearchService, pageSize int) ([]model.Product, int, error) {

	res, err := search.Do(context.Background())
	if err != nil {
		log.Fatal(err)
		return nil, 1, err
	}

	products, err := collectProducts(res.Hits)
	if err != nil {
		return products, 1, err
	}

	maxPage := res.TotalHits() / int64(pageSize)
	return products, int(maxPage + 1), nil
}

func collectProducts(hits *elastic.SearchHits) ([]model.Product, error) {

	products := make([]model.Product, 0)

	for _, hit := range hits.Hits {

		var p model.Product
		if err := json.Unmarshal(*hit.Source, &p); err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func collectAutocomplete(hits *elastic.SearchHits) ([]model.AutocompleteProduct, error) {

	products := make([]model.AutocompleteProduct, 0)

	for _, hit := range hits.Hits {

		var p model.AutocompleteProduct
		if err := json.Unmarshal(*hit.Source, &p); err != nil {
			return nil, err
		}

		highlight := hit.Highlight["pname"][0]
		p.Highlight = highlight

		products = append(products, p)
	}

	return products, nil
}