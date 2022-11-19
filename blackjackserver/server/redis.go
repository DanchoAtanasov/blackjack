package server

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/go-redis/redis"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var REDIS_HOST string = getEnv("REDIS_HOST", "localhost")
var REDIS_PORT string = getEnv("REDIS_PORT", "6379") // This about making int

type PlayerDetails struct {
	Name    string
	BuyIn   int
	CurrBet int
}

func fetchPlayerDetails(token string) PlayerDetails {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT),
		Password: "",
		DB:       0,
	})

	val, err := client.Get(token).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("got from redis %s\n", val)

	var pd PlayerDetails
	err = json.Unmarshal([]byte(val), &pd)
	if err != nil {
		fmt.Println("failed to unmarshal")
	}

	return pd
}

func getPlayerDetails(conn net.Conn) PlayerDetails {
	type Token struct {
		Token string
	}
	var token Token

	// Ask client for token
	msg := ReadData(conn)
	fmt.Printf("Received %s", msg)
	err := json.Unmarshal([]byte(msg), &token)
	if err != nil {
		fmt.Println("Failed to unmarshal")
		// TODO: handle failure
	}

	fmt.Println(token)
	return fetchPlayerDetails(token.Token)
}
