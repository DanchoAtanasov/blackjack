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

const FRONT_END_PORT int = 5173
const PORT int = 3333

var REDIS_HOST string = getEnv("REDIS_HOST", "localhost")
var PRIVATE_KEY *rsa.PrivateKey = loadPrivateKey()
var DOMAIN string = getEnv("DOMAIN", "blackjack.gg")
var FRONT_END_URL string = fmt.Sprintf("https://%s:%d", DOMAIN, FRONT_END_PORT)
var BLACKJACK_SERVER_PATH string = fmt.Sprintf("%s/blackjack/", DOMAIN)

var db UsersDatabase

type PlayerRequest struct {
	Name    string
	BuyIn   int
	CurrBet int
}

type LoginRequest struct {
	Username string
	Password string
}

type BlackjackServerDetails struct {
	GameServer string
	Token      string
}

func setCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", FRONT_END_URL)
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", "240")
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "Time to play some blackjack huh\n")
}

func new(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got %s /new request\n", r.Method)
	fmt.Println(r.Cookies())
	response, _ := json.Marshal("new")
	io.WriteString(w, string(response))
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got %s /login request\n", r.Method)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	fmt.Printf("body: %s\n", body)

	var loginRequest LoginRequest
	err = json.Unmarshal([]byte(body), &loginRequest)
	if err != nil {
		fmt.Printf("could not parse login request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("got %s, %s\n", loginRequest.Username, loginRequest.Password)
	// db.query()
	user, err := db.getUser(loginRequest.Username)

	if loginRequest.Password != user.Password {
		fmt.Println("Wrong password")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println("Authorized")

	redisToken := uuid.NewString()
	token := generateJwt(redisToken)

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(1 * time.Hour),
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
		Domain:   fmt.Sprintf(".%s", DOMAIN),
	})
}

func play(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got %s /play request\n", r.Method)
	fmt.Println(r.Cookies())

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

	redisToken := uuid.NewString()
	token := generateJwt(redisToken)
	fmt.Println(token)

	response := BlackjackServerDetails{GameServer: BLACKJACK_SERVER_PATH, Token: token}
	responseString, _ := json.Marshal(response)
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Hour),
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
		Domain:   fmt.Sprintf(".%s", DOMAIN),
	})
	fmt.Println("cookie set")
	io.WriteString(w, string(responseString))

	go storeSession(redisToken, playerRequest)

}

func generateJwt(uuid string) string {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"token": uuid,
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

// Middleware that handles CORS
type Cors struct {
	handler http.Handler
}

func (l *Cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setCorsHeaders(w)
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
	var err error
	db, err = NewUsersDatabase()
	if err != nil {
		fmt.Printf("error couldn't connect to database: %s\n", err)
		os.Exit(1)
	}
	defer db.Close()
	db.query()

	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/play", play)
	mux.HandleFunc("/new", new)
	mux.HandleFunc("/login", login)
	corsMux := NewCors(mux)

	fmt.Printf("Api server started on port %d\n", PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), corsMux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
