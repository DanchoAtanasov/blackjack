package server

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	settings "blackjack/configs"
)

func ReadData(conn net.Conn) string {
	msg, _, err := wsutil.ReadClientData(conn)
	if err != nil {
		fmt.Println("Read failed, ", err)
	}

	msg_str := string(msg)
	fmt.Println("Received: ", msg_str)
	return msg_str
}

func SendData(conn net.Conn, msg string) {
	err := wsutil.WriteServerMessage(conn, ws.OpText, []byte(msg))
	if err != nil {
		fmt.Println("Send failed, ", err)
	}
}

type Room struct {
	connections   []net.Conn
	currentPlayer int
	isFullCond    sync.Cond
	mutex         sync.Mutex
}

func MakeRoom() *Room {
	room := Room{}
	room.isFullCond = *sync.NewCond(&room.mutex)
	return &room
}

func (room *Room) GetCurrPlayerConn() *net.Conn {
	return &room.connections[room.currentPlayer]
}

func (room *Room) ChangePlayer() {
	room.currentPlayer += 1
	room.currentPlayer %= len(room.connections) // TODO if player disconnects
}

func (room *Room) SendAll(msg string) {
	for _, conn := range room.connections {
		SendData(conn, msg)
	}
}

func (room *Room) WaitForPlayers() {
	fmt.Println("Waiting for players")
	room.isFullCond.L.Lock()
	for len(room.connections) != settings.RoomSize {
		room.isFullCond.Wait()
	}

	room.isFullCond.L.Unlock()
	fmt.Println("Wait is over!")
}

type Server struct {
	newConnectionMutex sync.Mutex
	room               *Room
}

func (server *Server) registerPlayer(conn *net.Conn) {
	fmt.Println("Registering player")
	server.newConnectionMutex.Lock()
	defer server.newConnectionMutex.Unlock()

	currRoom := server.room
	currRoom.connections = append(currRoom.connections, *conn)

	if len(currRoom.connections) == settings.RoomSize {
		currRoom.isFullCond.Broadcast()
		server.room = MakeRoom()
	}
}

func MakeServer() *Server {
	server := Server{}
	server.room = MakeRoom()
	return &server
}

func (server *Server) GetRoom() *Room {
	return server.room
}

func (server *Server) Serve() {
	port := 8080
	fmt.Println("Serving on port", port)

	http.ListenAndServe(fmt.Sprintf(":%d", port),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, _, _, err := ws.UpgradeHTTP(r, w)
			if err != nil {
				// handle error
			}

			fmt.Println("New player connected")
			go server.registerPlayer(&conn)
		}))
}
