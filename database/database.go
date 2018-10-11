package database

import (
	"github.com/olivere/elastic"
	"gopkg.in/mgo.v2"
)

const pageSize = 24

type Database struct {
	products      *mgo.Collection
	users         *mgo.Collection
	carts         *mgo.Collection
	orders        *mgo.Collection
	elasticClient *elastic.Client
}

func New(mongoSession *mgo.Session, elasticClient *elastic.Client) Database {

	mongodb := mongoSession.DB("shop")

	products := mongodb.C("products")
	users := mongodb.C("users")
	carts := mongodb.C("carts")
	orders := mongodb.C("orders")

	return Database{
		products: products,
		users: users,
		carts: carts,
		orders: orders,
		elasticClient: elasticClient,
	}
}

