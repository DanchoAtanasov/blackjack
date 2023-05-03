package routes

import (
	"apiserver/tokens"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var UserTokenNotFound = errors.New("user-token not found")
var UserTokenCannotBeParsed = errors.New("user-token cannot be parsed")

type CookieLoginResponse struct {
	Userid   string
	Username string
}

func parseUserTokenFromRequest(r *http.Request) (tokens.UserToken, error) {
	cookie, err := r.Cookie("user-token")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Println("user-token cookie not found")
			return tokens.UserToken{}, UserTokenNotFound
		} else {
			fmt.Println(err)
			return tokens.UserToken{}, UserTokenNotFound
		}
	}

	userToken, err := tokens.ParseUserToken(cookie.Value)
	if err != nil {
		return tokens.UserToken{}, UserTokenCannotBeParsed
	}
	return userToken, nil
}

func (h *RouteHandler) CookieLogin(w http.ResponseWriter, r *http.Request) {
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

	resp := CookieLoginResponse{Userid: userToken.UserId, Username: userToken.Username}
	response, _ := json.Marshal(resp)
	io.WriteString(w, string(response))
}
