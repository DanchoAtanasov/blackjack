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

type PlayerRequest struct {
	Name  string
	BuyIn int
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "Time to play some blackjack huh\n")
}

func play(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /play request\n")

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
	type PlayerRequestResponse struct {
		GameServer string
		Token      string
	}
	response := PlayerRequestResponse{GameServer: "localhost:8080", Token: uuid.NewString()}
	responseString, _ := json.Marshal(response)
	io.WriteString(w, string(responseString))

	go storeSession(response.Token, playerRequest)
}

func storeSession(token string, playerRequest PlayerRequest) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
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
	val, err := client.Get(token).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("got from redis %s\n", val)

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/play", play)

	err := http.ListenAndServe(":3333", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
