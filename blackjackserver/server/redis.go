package server

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v5"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func loadPublicKey() *rsa.PublicKey {
	keyData, _ := os.ReadFile("keys/key.pub")
	key, _ := jwt.ParseRSAPublicKeyFromPEM(keyData)
	return key
}

var PUBLIC_KEY *rsa.PublicKey = loadPublicKey()

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

func parseJwt(tokenString string) string {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return PUBLIC_KEY, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["token"])
		return claims["token"].(string)
	} else {
		fmt.Println(err)
		return ""
	}
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
	redisToken := parseJwt(token.Token)
	return fetchPlayerDetails(redisToken)
}
