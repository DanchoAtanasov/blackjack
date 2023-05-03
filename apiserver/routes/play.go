package routes

import (
	env "apiserver/environment"
	"apiserver/tokens"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

type PlayerRequest struct {
	CurrBet int
}

type BlackjackServerDetails struct {
	GameServer string
	Token      string
}

type PlayerSessionInformation struct {
	Name    string
	BuyIn   int
	CurrBet int
}

func (h *RouteHandler) Play(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println("User token is valid")

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
	fmt.Printf("got %s, %d\n", userToken.Username, playerRequest.CurrBet)

	fmt.Println("Getting user from database")
	user, err := h.db.GetUser(userToken.Username)
	if err != nil {
		fmt.Printf("Cannot get user %s, %s\n", userToken.Username, err)
		return
	}

	redisToken := uuid.NewString()
	sessionToken := tokens.GenerateSessionToken(redisToken)

	http.SetCookie(w, &http.Cookie{
		Name:     "session-token",
		Value:    sessionToken,
		Path:     "/",
		Expires:  time.Now().Add(1 * time.Hour),
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
		Domain:   fmt.Sprintf(".%s", env.DOMAIN),
	})

	response := BlackjackServerDetails{GameServer: env.BLACKJACK_SERVER_PATH, Token: sessionToken}
	responseString, _ := json.Marshal(response)
	io.WriteString(w, string(responseString))
	playerSessionInformation := PlayerSessionInformation{
		Name:    userToken.Username,
		BuyIn:   user.BuyIn,
		CurrBet: playerRequest.CurrBet,
	}

	go storeSession(redisToken, playerSessionInformation)

}

func storeSession(token string, playerSessionInformation PlayerSessionInformation) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", env.REDIS_HOST),
		Password: "",
		DB:       0,
	})

	json, err := json.Marshal(playerSessionInformation)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.Set(token, json, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}
