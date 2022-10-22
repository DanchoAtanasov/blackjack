package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var REDIS_HOST string = getEnv("REDIS_HOST", "localhost")

type PlayerRequest struct {
	Name  string
	BuyIn int
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "Time to play some blackjack huh\n")
}

func play(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got %s /play request\n", r.Method)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", "240")

	if r.Method == "OPTIONS" {
		io.WriteString(w, "")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body %s\n", err)
	}
	fmt.Printf("body: %s\n", body)

	var playerRequest PlayerRequest
	err = json.Unmarshal([]byte(body), &playerRequest)
	if err != nil {
		fmt.Printf("could not parse body")
		return
	}
	fmt.Printf("got %s, %d\n", playerRequest.Name, playerRequest.BuyIn)
	type BlackjackServerDetails struct {
		GameServer string
		Token      string
	}
	response := BlackjackServerDetails{GameServer: "localhost:8080", Token: uuid.NewString()}
	responseString, _ := json.Marshal(response)
	io.WriteString(w, string(responseString))

	go storeSession(response.Token, playerRequest)
}

func storeSession(token string, playerRequest PlayerRequest) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", REDIS_HOST),
		Password: "",
		DB:       0,
	})

	json, err := json.Marshal(playerRequest)
	if err != nil {
		fmt.Println(err)
	}

	err = client.Set(token, json, 0).Err()
	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	port := 3333
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/play", play)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
