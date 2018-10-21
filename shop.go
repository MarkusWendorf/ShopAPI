package main

import (
	"flag"
	"github.com/go-chi/chi"
	"github.com/olivere/elastic"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"shopApi/database"
	"shopApi/middleware"
	"shopApi/routes"
	"time"
)

// time to wait for containers to spin up
const StartUpTimout = 60 * time.Second

func main() {

	mongo := "mongo:27017"
	elasticSearch := "http://elastic:9200"

	runSetup := flag.Bool("setup-only", false, "only run intial setup")

	flag.Parse()
	if *runSetup {
		log.Println("Running initial setup")
		setup(mongo, elasticSearch)
	} else {
		serve(mongo, elasticSearch,"3000")
	}
}

func serve(mongoAddr string, elasticAddr string, apiPort string) {

	mongoSession := initMongoDB(mongoAddr)
	defer mongoSession.Close()

	elasticClient := initElasticSearch(elasticAddr)

	db := database.New(mongoSession, elasticClient)
	handlers := routes.NewHandlers(db)

	mux := chi.NewMux()

	corsOptions := middleware.CorsConfig{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   "authorization",
	}

	mux.Use(middleware.DelayMiddleware)
	mux.Use(middleware.CorsMiddleware(corsOptions))

	mux.Get("/products/{id}", handlers.GetProduct)
	mux.Get("/autocomplete/{search}", handlers.Autocomplete)
	mux.Get("/categories", handlers.GetCategoryNames)
	mux.Get("/search", handlers.QueryProducts)
	mux.Post("/login", handlers.Login)
	mux.Post("/register", handlers.Register)

	// auth required
	mux.Group(func(r chi.Router) {

		r.Use(middleware.AuthMiddleware)

		r.Post("/cart", handlers.Cart)
		r.Post("/order", handlers.Order)
	})

	log.Printf("Shop API running on port %s\n", apiPort)
	log.Fatal(http.ListenAndServe(":"+apiPort, mux))
}

func initElasticSearch(url string) *elastic.Client {

	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetHealthcheckTimeoutStartup(StartUpTimout))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB:", url)
	return client
}

func initMongoDB(url string) *mgo.Session {

	session, err := mgo.DialWithTimeout(url, StartUpTimout)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to ElasticSearch:", url)
	return session
}
