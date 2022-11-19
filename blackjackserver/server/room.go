package server

import (
	settings "blackjack/config"
	"blackjack/models"
	"fmt"
	"net"
	"os"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Struct used to couple the player model and the connection to them
type PlayerConn struct {
	player models.Player
	Conn   net.Conn
}

type Room struct {
	Log         *logrus.Logger
	Id          string
	playerConns []PlayerConn
}

func MakeRoom() *Room {
	room := Room{}
	room.Id = uuid.New().String()
	room.Log = MakeLog(room.Id)
	room.playerConns = make([]PlayerConn, 0, settings.RoomSize)

	return &room
}

func (room *Room) GetCurrPlayerConn() *PlayerConn {
	return &room.playerConns[0]
}

func (room *Room) RemoveDisconnectedPlayer() {
	// TODO improve player diconnecting/rotating logic
	// maybe couple the player and room objects better
	room.playerConns = room.playerConns[1:]
}

func (room *Room) ChangePlayer() {
	// Change players by popping from queue and appending
	var currentConn PlayerConn
	currentConn, room.playerConns = room.playerConns[0], room.playerConns[1:]
	room.playerConns = append(room.playerConns, currentConn)
}

func (room *Room) SendAll(msg string) {
	for _, playerConnection := range room.playerConns {
		SendData(playerConnection.Conn, msg)
	}
}

func (room *Room) GetPlayers() []models.Player {
	var players []models.Player
	for i := 0; i < len(room.playerConns); i++ {
		players = append(players, room.playerConns[i].player)
	}
	return players
}

func MakeLog(id string) *logrus.Logger {
	filename := fmt.Sprintf("./logs/%s.log", id)
	log := logrus.New()
	log.Out = os.Stdout

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	return log
}
