package routes

import (
	env "apiserver/environment"
	hashing "apiserver/hashing"
	tokens "apiserver/tokens"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type LoginRequest struct {
	Username string
	Password string
}

func (h *RouteHandler) Login(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.db.GetUser(loginRequest.Username)
	if err != nil {
		fmt.Printf("Cannot get user %s, %s\n", loginRequest.Username, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !hashing.CheckPasswordHash(loginRequest.Password, user.Password) {
		fmt.Println("Wrong password")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println("Authorized")

	userToken := tokens.GenerateUserToken(user.Id, user.Username)

	http.SetCookie(w, &http.Cookie{
		Name:     "user-token",
		Value:    userToken,
		Path:     "/",
		Expires:  time.Now().Add(1 * time.Hour),
		Secure:   true,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
		Domain:   fmt.Sprintf(".%s", env.DOMAIN),
	})
}
