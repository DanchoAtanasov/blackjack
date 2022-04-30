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

type Room struct {
	connections   []net.Conn
	currentPlayer int
	isFullCond    sync.Cond
	mutex         sync.Mutex
}

func makeRoom() Room {
	room := Room{}
	room.isFullCond = *sync.NewCond(&room.mutex)
	return room // TODO could copy lock
}

func (room *Room) GetCurrPlayerConn() *net.Conn {
	return &room.connections[room.currentPlayer]
}

func (room *Room) ChangePlayer() {
	room.currentPlayer += 1
	room.currentPlayer %= len(room.connections)
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
	rooms              []Room
}

func (server *Server) registerPlayer(conn *net.Conn) {
	server.newConnectionMutex.Lock()
	lastRoomIdx := len(server.rooms) - 1
	currRoom := &server.rooms[lastRoomIdx]
	currRoom.connections = append(currRoom.connections, *conn)
	fmt.Println(len(currRoom.connections))
	if len(currRoom.connections) == settings.RoomSize {
		server.rooms = append(server.rooms, makeRoom())
		fmt.Println("Broadcasting")
		currRoom.isFullCond.Broadcast()
	}
	server.newConnectionMutex.Unlock()
}

func MakeServer() Server {
	server := Server{}
	server.rooms = append(server.rooms, makeRoom())
	return server // TODO: fix this copy of lock
}

func (server *Server) GetLastRoom() *Room {
	return &server.rooms[len(server.rooms)-1]
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
