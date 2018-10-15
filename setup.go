package main

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/olivere/elastic"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"log"
	"os"
	"shopApi/model"
	"strconv"
)

func setup(mongodb string, elasticSearch string) {

	products := parseProductsJSON()
	log.Println("PRODUCTS:::::", len(products))


	setupMongo(mongodb, products)
	setupElastic(elasticSearch, products)
}

func setupMongo(url string, products []model.Product) {

	log.Println("Setup MongoDB", url)

	session := initMongoDB(url)
	defer session.Close()

	db := session.DB("shop")
	productsCollection := db.C("products")

	for _, p := range products {

		err := productsCollection.Insert(p)
		if err != nil {
			log.Fatal(err)
		}
	}

	index := mgo.Index{Key: []string{"id"}, Unique: true, DropDups: true, Sparse: true,}
	if err := productsCollection.EnsureIndex(index); err != nil {
		log.Fatal(err)
	}
}

func setupElastic(url string, products []model.Product) {

	log.Println("Setup Elasticsearch")

	client := initElasticSearch(url)

	index := "products"
	b, err := ioutil.ReadFile("mapping.json")

	resp, err := client.CreateIndex(index).BodyString(string(b)).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if resp.Acknowledged {
		log.Println("Index successfully created")
	} else {
		log.Fatal("Index creation failed")
	}

	bulk := client.Bulk()
	for i, p := range products {
		req := elastic.NewBulkIndexRequest().
			Index(index).
			Type("product").
			Id(strconv.Itoa(i)).
			Doc(p)

		bulk.Add(req)
	}

	result, err := bulk.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if result.Errors {
		log.Fatal(result.Failed()[0].Error)
	}

	log.Println("Indexing completed")
}

func parseProductsJSON() []model.Product {

	f, err := os.Open("products.json")
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(f)

	var products struct {
		Products []model.Product
	}

	err = json.NewDecoder(reader).Decode(&products)
	if err != nil {
		panic(err)
	}

	ps := products.Products
	for i := 0; i < len(ps); i++ {
		ps[i].ImgUrl = "/images/" + ps[i].Id + ".jpg"
	}

	return products.Products
}