package server

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	settings "blackjack/configs"
)

func ReadData(conn net.Conn) string {
	msg, _, err := wsutil.ReadClientData(conn)
	if err != nil {
		// handle error
	}
	msg_str := string(msg)
	fmt.Println("Received: ", msg_str)
	return msg_str
}

func SendData(conn net.Conn, msg string) {
	err := wsutil.WriteServerMessage(conn, ws.OpText, []byte(msg))
	if err != nil {
		// handle error
	}
}

type Server struct {
	connections        []net.Conn
	newConnectionMutex sync.Mutex
	someCond           sync.Cond
	currentPlayer      int
}

func (server *Server) GetCurrPlayerConn() *net.Conn {
	currPlayerConn := server.connections[server.currentPlayer]
	return &currPlayerConn
}

func (server *Server) ChangePlayer() {
	server.currentPlayer += 1
	server.currentPlayer %= len(server.connections)
}

func (server *Server) SendAll(msg string) {
	for _, conn := range server.connections {
		SendData(conn, msg)
	}
}

func (server *Server) registerPlayer(conn *net.Conn) {
	server.newConnectionMutex.Lock()
	server.connections = append(server.connections, *conn)
	if len(server.connections) == settings.RoomSize {
		server.someCond.Broadcast()
	}
	server.newConnectionMutex.Unlock()
}

func (server *Server) WaitForPlayers() {
	fmt.Println("Waiting for players")
	server.someCond.L.Lock()
	for len(server.connections) != settings.RoomSize {
		server.someCond.Wait()
	}
	server.someCond.L.Unlock()
	fmt.Println("Wait is over!")
}

func NewServer() Server {
	server := Server{}
	server.someCond = *sync.NewCond(&server.newConnectionMutex)
	return server // TODO: fix this copy of lock
}

func (server *Server) Serve() {
	port := 8080
	port_str := ":" + strconv.Itoa(port)
	fmt.Println("Serving on port", port)

	http.ListenAndServe(port_str, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			// handle error
		}

		fmt.Println("New player connected")
		go server.registerPlayer(&conn)
	}))
}
