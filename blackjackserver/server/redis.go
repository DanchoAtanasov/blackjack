package server

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
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
var REDIS_PORT string = getEnv("REDIS_PORT", "6379")

type Token struct {
	Token string
}

type PlayerDetails struct {
	Name    string
	BuyIn   int
	CurrBet int
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

func getPlayerDetails(sessionId string) PlayerDetails {
	// TODO: reuse an existing redis client connection
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT),
		Password: "",
		DB:       0,
	})
	defer client.Close()

	val, err := client.Get(sessionId).Result()
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

func SetPlayerDetails(sessionId string, pd PlayerDetails) {
	fmt.Printf("Setting player details for session: %s\n", sessionId)
	// TODO: reuse an existing redis client connection
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT),
		Password: "",
		DB:       0,
	})
	defer client.Close()

	pdJson, err := json.Marshal(pd)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.Set(sessionId, pdJson, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}
