package server

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type NumPlayers struct {
	mutex sync.Mutex
	value int
}

func gameplayLoop() {

}

func readData(conn net.Conn) []byte {
	msg, _, err := wsutil.ReadClientData(conn)
	if err != nil {
		// handle error
	}
	fmt.Println("Received: ", msg)
	return msg
}

func sendData(conn net.Conn, msg []byte) {
	err := wsutil.WriteServerMessage(conn, ws.OpBinary, msg)
	if err != nil {
		// handle error
	}
}

func findPlayers(conn net.Conn, numPlayers *NumPlayers, cond *sync.Cond) {
	fmt.Println("In find")
	_ = readData(conn)
	var resp []byte

	numPlayers.mutex.Lock()
	{
		if numPlayers.value < 2 {
			numPlayers.value += 1
			fmt.Println("Players are now ", numPlayers.value)
			resp = []byte("OK")
		} else {
			resp = []byte("FULL")
			fmt.Println("Room is full ", numPlayers.value)
		}
	}
	if numPlayers.value == 2 {
		cond.Broadcast()
	}
	numPlayers.mutex.Unlock()

	sendData(conn, resp)

	cond.L.Lock()
	for numPlayers.value < 2 {
		cond.Wait()
	}
	fmt.Println("Got 2, lets' play")
	cond.L.Unlock()

	resp = []byte("Im going to start dealing")
	sendData(conn, resp)
}

type Server struct {
	connections []net.Conn
}

func (server *Server) Serve() {
	port := 8080
	fmt.Println("Serving on port", port)

	numPlayers := NumPlayers{value: 0}
	cond := sync.NewCond(&numPlayers.mutex)
	http.ListenAndServe(":"+strconv.Itoa(port), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			// handle error
		}
		fmt.Println("New player connected")
		// go doNothing(conn, &numPlayers)
		go findPlayers(conn, &numPlayers, cond)
	}))
}
