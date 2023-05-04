package routes

import (
	"apiserver/database"
	"apiserver/session_cache"
)

// Handler that couples the path handling with the database connection and session cache
type RouteHandler struct {
	db *database.UsersDatabase
	sc *sessioncache.SessionCache
}

func NewRouteHandler(db *database.UsersDatabase, sc *sessioncache.SessionCache) *RouteHandler {
	return &RouteHandler{
		db: db,
		sc: sc,
	}
}
