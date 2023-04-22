package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Same as LoginRequest for now, could add password reenter
type SignupRequest struct {
	Username string
	Password string
}

func (h *RouteHandler) Signup(w http.ResponseWriter, r *http.Request) {
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

	_, err = h.db.AddUser(signupRequest.Username, signupRequest.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
