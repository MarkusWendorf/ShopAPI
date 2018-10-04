package routes

import "shopApi/database"

type Handlers struct {
	db database.Database
}

func NewHandlers(db database.Database) Handlers {
	return Handlers{db: db}
}
