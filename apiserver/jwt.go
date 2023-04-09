package main

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func loadPublicKey() *rsa.PublicKey {
	keyData, _ := os.ReadFile("keys/key.pub")
	key, _ := jwt.ParseRSAPublicKeyFromPEM(keyData)
	return key
}

var PUBLIC_KEY *rsa.PublicKey = loadPublicKey()

func generateUserToken(userId string, username string) string {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"userid":   userId,
		"username": username,
		"nbf":      time.Now().Unix(),
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(PRIVATE_KEY)
	if err != nil {
		fmt.Printf("Token signing failed %s/n", err)
		return ""
	}
	return tokenString
}

func generateSessionToken(redisToken string) string {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"token": redisToken,
		"nbf":   time.Now().Unix(),
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(1 * time.Hour).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(PRIVATE_KEY)
	if err != nil {
		fmt.Printf("Token signing failed %s/n", err)
		return ""
	}
	return tokenString
}

func parseUserToken(tokenString string) (UserToken, error) {
	claims, err := parseJwt(tokenString)
	if err != nil {
		fmt.Println("User token invalid")
		return UserToken{}, err
	}
	return UserToken{UserId: claims["userid"].(string), Username: claims["username"].(string)}, nil
}

func parseSessionToken(tokenString string) (SessionToken, error) {
	claims, err := parseJwt(tokenString)
	if err != nil {
		fmt.Println("Session token invalid")
		return SessionToken{}, err
	}
	return SessionToken{Token: claims["token"].(string)}, nil
}

func parseJwt(tokenString string) (jwt.MapClaims, error) {
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
		return claims, nil
	} else {
		fmt.Println(err)
		return claims, errors.New("Invalid jwt")
	}
}
