package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type EndSessionRequest struct {
	SessionId string
}

func (h *RouteHandler) EndSession(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not read body %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	var endSessionRequest EndSessionRequest
	err = json.Unmarshal([]byte(body), &endSessionRequest)
	if err != nil {
		fmt.Printf("could not parse end session request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	playerSession := h.sc.GetSession(endSessionRequest.SessionId)

	err = h.db.SetUserBuyIn(playerSession.Name, playerSession.BuyIn)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Buyin set")

	fmt.Println("Deleting session")
	go h.sc.DeleteSession(endSessionRequest.SessionId)
}
