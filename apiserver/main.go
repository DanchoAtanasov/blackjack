package main

import (
	"apiserver/database"
	env "apiserver/environment"
	"apiserver/routes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func setCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", env.FRONT_END_URL)
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", "240")
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Time to play some blackjack huh\n")
}

// Middleware that handles CORS
type Cors struct {
	handler http.Handler
}

func (l *Cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setCorsHeaders(w)
	fmt.Printf("Got %s %s\n", r.Method, r.RequestURI)
	if r.Method == "OPTIONS" {
		io.WriteString(w, "")
		return
	}
	l.handler.ServeHTTP(w, r)
}

func NewCors(handlerToWrap http.Handler) *Cors {
	return &Cors{handlerToWrap}
}

func main() {
	db, err := database.NewUsersDatabase()
	if err != nil {
		fmt.Printf("error couldn't connect to database: %s\n", err)
		os.Exit(1)
	}
	defer db.Close()

	routeHandler := routes.NewRouteHandler(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/play", routeHandler.Play)
	mux.HandleFunc("/cookie", routeHandler.CookieLogin)
	mux.HandleFunc("/login", routeHandler.Login)
	mux.HandleFunc("/signup", routeHandler.Signup)
	mux.HandleFunc("/endsession", routeHandler.EndSession)
	corsMux := NewCors(mux)

	fmt.Printf("Api server started on port %d\n", env.PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%d", env.PORT), corsMux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
