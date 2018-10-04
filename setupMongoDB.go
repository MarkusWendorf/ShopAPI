package main

import (
	"shopApi/model"
	"bufio"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
)

func setup(host string) {

	fmt.Printf("Running initial mongodb setup %s\n", host)

	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	db := session.DB("shop")

	productsCollection := db.C("products")

	products := parseProductsJSON()
	fmt.Printf("Found %d products", len(products))

	for _, p := range products {

		p.ImgUrl = "/images/" + p.Id + ".jpg"

		err = productsCollection.Insert(p)
		if err != nil {
			panic(err)
		}

	}

	index := mgo.Index{
		Key:      []string{"id"},
		Unique:   true,
		DropDups: true,
		Sparse:   true,
	}

	err = productsCollection.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	fmt.Println("Products collection is populated\nSetup completed")
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

	return products.Products
}
