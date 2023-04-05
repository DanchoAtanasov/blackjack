package main

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func loadPrivateKey() *rsa.PrivateKey {
	keyData, _ := os.ReadFile("keys/key.pem")
	key, _ := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	return key
}

var REDIS_HOST string = getEnv("REDIS_HOST", "localhost")
var PRIVATE_KEY *rsa.PrivateKey = loadPrivateKey()

type PlayerRequest struct {
	Name    string
	BuyIn   int
	CurrBet int
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
	fmt.Printf("got %s, %d, %d\n", playerRequest.Name, playerRequest.BuyIn, playerRequest.CurrBet)

	type BlackjackServerDetails struct {
		GameServer string
		Token      string
	}

	redisToken := uuid.NewString()
	token := generateJwt(playerRequest, redisToken)
	fmt.Println(token)

	response := BlackjackServerDetails{GameServer: "localhost/blackjack/", Token: token}
	responseString, _ := json.Marshal(response)
	io.WriteString(w, string(responseString))

	go storeSession(redisToken, playerRequest)

}

func generateJwt(playerRequest PlayerRequest, uuid string) string {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"token": uuid,
		// "name": playerRequest.Name,
		// "buyin":   playerRequest.BuyIn,
		// "currbet": playerRequest.CurrBet,
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(PRIVATE_KEY)
	if err != nil {
		fmt.Printf("Token signing failed %s/n", err)
		return ""
	}

	fmt.Printf("Generated token: %s\n", tokenString)
	return tokenString
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

	fmt.Printf("Api server started on port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
