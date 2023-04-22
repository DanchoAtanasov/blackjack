package routes

import "apiserver/database"

// Handler that couples the path handling with the database connection
type RouteHandler struct {
	db *database.UsersDatabase
}

func NewRouteHandler(db *database.UsersDatabase) *RouteHandler {
	return &RouteHandler{
		db: db,
	}
}
