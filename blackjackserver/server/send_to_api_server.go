package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const API_SERVER_URL = "http://apiserver:3333/endsession"

type EndSessionRequest struct {
	SessionId string `json:"sessionid"`
}

func SendEndSession(sessionId string) {
	endSessionRequest := EndSessionRequest{SessionId: sessionId}
	fmt.Println(endSessionRequest)
	msg, err := json.Marshal(endSessionRequest)
	if err != nil {
		fmt.Printf("Can't marshal end session request for session id: %s", sessionId)
		return
	}

	req, err := http.NewRequest("POST", API_SERVER_URL, bytes.NewBuffer(msg))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
