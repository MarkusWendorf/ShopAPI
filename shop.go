package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"shopApi/database"
	"shopApi/middleware"
	"shopApi/routes"
)

func main() {

	runSetup := flag.Bool("setup-only", false, "run initial setup (only run once)")
	ipFlag := flag.String("ip", "localhost", "")
	apiPort := flag.String("port", "3000", "")
	flag.Parse()

	ip := *ipFlag
	port := *apiPort
	if *runSetup {
		setup(ip)
		return
	}

	serve(ip, port)
}

func serve(host string, port string) {

	fmt.Printf("Connecting to mongodb %s\n", host)

	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	mongodb := session.DB("shop")
	products := mongodb.C("products")
	users := mongodb.C("users")
	carts := mongodb.C("carts")
	orders := mongodb.C("orders")

	db := database.New(products, users, carts, orders)
	handlers := routes.NewHandlers(db)

	mux := chi.NewMux()
	mux.Use(middleware.CorsMiddleware)

	mux.Get("/products/{id}", handlers.GetProduct)
	mux.Get("/categories", handlers.GetCategoryNames)
	mux.Get("/query", handlers.QueryProducts)
	mux.Post("/login", handlers.Login)
	mux.Post("/register", handlers.Register)

	// auth required
	mux.Group(func(r chi.Router) {

		r.Use(middleware.AuthMiddleware)

		r.Post("/cart", handlers.Cart)
		r.Post("/order", handlers.Order)
	})

	fmt.Printf("Shop API running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
