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
	// Name    string
	BuyIn   int
	CurrBet int
}

type PlayerSessionInformation struct {
	Name string
	PlayerRequest
}

type LoginRequest struct {
	Username string
	Password string
}

// Same as LoginRequest for now, could add password reenter
type SignupRequest struct {
	Username string
	Password string
}

type BlackjackServerDetails struct {
	GameServer string
	Token      string
}

type CookieLoginResponse struct {
	// Token    string
	Useid    string
	Username string
}

var UserTokenNotFound = errors.New("user-token not found")
var UserTokenCannotBeParsed = errors.New("user-token cannot be parsed")

func setCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", FRONT_END_URL)
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", "240")
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Time to play some blackjack huh\n")
}

func parseUserTokenFromRequest(r *http.Request) (UserToken, error) {
	cookie, err := r.Cookie("user-token")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Println("user-token cookie not found")
			return UserToken{}, UserTokenNotFound
		} else {
			fmt.Println(err)
			return UserToken{}, UserTokenNotFound
		}
	}

	userToken, err := parseUserToken(cookie.Value)
	if err != nil {
		return UserToken{}, UserTokenCannotBeParsed
	}
	return userToken, nil
}

func cookieLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Cookies())
	userToken, err := parseUserTokenFromRequest(r)
	if err != nil {
		fmt.Println("cookie not valid")
		if err == UserTokenNotFound {
			// TODO: revise eror codes
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	resp := CookieLoginResponse{Useid: userToken.UserId, Username: userToken.Username}
	response, _ := json.Marshal(resp)
	io.WriteString(w, string(response))
}

func login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	var loginRequest LoginRequest
	err = json.Unmarshal([]byte(body), &loginRequest)
	if err != nil {
		fmt.Printf("could not parse login request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := db.getUser(loginRequest.Username)
	if err != nil {
		fmt.Printf("Cannot get user %s, %s\n", loginRequest.Username, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !CheckPasswordHash(loginRequest.Password, user.Password) {
		fmt.Println("Wrong password")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println("Authorized")

	userToken := generateUserToken(user.Id, user.Username)

	http.SetCookie(w, &http.Cookie{
		Name:     "user-token",
		Value:    userToken,
		Path:     "/",
		Expires:  time.Now().Add(1 * time.Hour),
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
		Domain:   fmt.Sprintf(".%s", DOMAIN),
	})
}

func signup(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	var signupRequest SignupRequest
	err = json.Unmarshal([]byte(body), &signupRequest)
	if err != nil {
		fmt.Printf("could not parse singup request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.addUser(signupRequest.Username, signupRequest.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func play(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Cookies())

	userToken, err := parseUserTokenFromRequest(r)
	if err != nil {
		fmt.Println("cookie not valid")
		if err == UserTokenNotFound {
			// TODO: revise eror codes
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	fmt.Println("User token in valid")

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
	fmt.Printf("got %s, %d, %d\n", userToken.Username, playerRequest.BuyIn, playerRequest.CurrBet)

	redisToken := uuid.NewString()
	sessionToken := generateSessionToken(redisToken)

	http.SetCookie(w, &http.Cookie{
		Name:     "session-token",
		Value:    sessionToken,
		Path:     "/",
		Expires:  time.Now().Add(1 * time.Hour),
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
		Domain:   fmt.Sprintf(".%s", DOMAIN),
	})

	response := BlackjackServerDetails{GameServer: BLACKJACK_SERVER_PATH, Token: sessionToken}
	responseString, _ := json.Marshal(response)
	io.WriteString(w, string(responseString))
	playerSessionInformation := PlayerSessionInformation{
		Name:          userToken.Username,
		PlayerRequest: playerRequest,
	}

	go storeSession(redisToken, playerSessionInformation)

}

func storeSession(token string, playerSessionInformation PlayerSessionInformation) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", REDIS_HOST),
		Password: "",
		DB:       0,
	})

	json, err := json.Marshal(playerSessionInformation)
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
	mux.HandleFunc("/cookie", cookieLogin)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/signup", signup)
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
