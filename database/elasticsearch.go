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

func getResults(search *elastic.SearchService, pageSize int) ([]model.Product, int, error) {

	res, err := search.Do(context.Background())
	if err != nil {
		log.Fatal(err)
		return nil, 1, err
	}

	products := make([]model.Product, 0)

	for _, hit := range res.Hits.Hits {

		var p model.Product
		if err = json.Unmarshal(*hit.Source, &p); err != nil {
			return nil, 1, err
		}

		products = append(products, p)

	}

	maxPage := res.TotalHits() / int64(pageSize)
	return products, int(maxPage + 1), nil
}

func (db *Database) ExecuteQuery(builder *QueryBuilder, pageSize int, page int) ([]model.Product, int, error) {
	return builder.Execute(db.elasticClient, pageSize, page)
}
