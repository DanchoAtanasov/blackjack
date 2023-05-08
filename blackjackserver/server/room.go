package server

import (
	settings "blackjack/config"
	"blackjack/messages"
	"blackjack/models"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Struct used to couple the player model and the connection to them
type PlayerConn struct {
	sessionId string
	player    *models.Player
	Conn      net.Conn
}

func (playerConn PlayerConn) saveDisconnectedPlayerDetails() {
	player := playerConn.player
	pd := PlayerDetails{
		// Don't know why implicit casting doesn't work
		Name:    player.Name,
		BuyIn:   player.BuyIn,
		CurrBet: player.CurrBet,
	}
	fmt.Println("Updating player details for disconnected player")
	setPlayerDetails(playerConn.sessionId, pd)
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

func (room *Room) IsEmpty() bool {
	return len(room.playerConns) == 0
}

func (room *Room) GetCurrPlayerConn() *PlayerConn {
	return &room.playerConns[0]
}

func (room *Room) AddNewPlayerConn(newPlayerConn PlayerConn) {
	room.playerConns = append(room.playerConns, newPlayerConn)
}

func (room *Room) RemoveDisconnectedPlayer(position int) {
	room.playerConns[position].saveDisconnectedPlayerDetails()
	SendEndSession(room.playerConns[position].sessionId)

	// TODO improve player diconnecting/rotating logic
	if len(room.playerConns) <= 1 {
		room.playerConns = room.playerConns[:0]
		return
	}
	// TODO: improve this as it's inefficient
	room.playerConns = append(room.playerConns[:position], room.playerConns[position+1:]...)
}

func (room *Room) ChangePlayer() {
	// Change players by popping from queue and appending
	// TODO: fix if player disconnects the order of players is offset
	// i.e. [1, 2, 3] -> 3 disconnects -> [1, 2] -> change player -> start turn [2, 1] instead of
	// [1, 2]
	if room.IsEmpty() {
		return
	}
	var currentConn PlayerConn
	currentConn, room.playerConns = room.playerConns[0], room.playerConns[1:]
	room.playerConns = append(room.playerConns, currentConn)
}

func (room *Room) SendAll(msg string) {
	for i := range room.playerConns {
		SendData(room.playerConns[i].Conn, msg)
	}
}

func (room *Room) ReadInMessages() {
	for i := range room.playerConns {
		currPlayer := room.playerConns[i].player

		// TODO add retry
		message := ReadData(room.playerConns[i].Conn)
		if message == "EOF" {
			fmt.Println("Player has disconnected")
			currPlayer.Active = false
			room.RemoveDisconnectedPlayer(i)
			continue
		}

		response, err := messages.DecodePlayerInMessage(message)
		if err != nil {
			fmt.Printf("Wrong player in response msg: %e\n", err)
		}

		if response.Playing {
			currPlayer.Active = true
			currPlayer.CurrBet = response.CurrentBet
		} else {
			currPlayer.Active = false
		}
	}
}

func (room *Room) SendCurrPlayer(msg string) {
	SendData(room.GetCurrPlayerConn().Conn, msg)
}

func (room *Room) ReadCurrPlayer() string {
	return ReadData(room.GetCurrPlayerConn().Conn)
}

func (room *Room) GetPlayers() []*models.Player {
	var players []*models.Player
	for i := 0; i < len(room.playerConns); i++ {
		players = append(players, room.playerConns[i].player)
	}
	return players
}

func MakeLog(id string) *logrus.Logger {
	filename := fmt.Sprintf("./logs/%s.log", id)
	log := logrus.New()

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		// Log to file and Stdout in the same time for easy debugging
		mw := io.MultiWriter(os.Stdout, file)

		// Set the logger's output to the multiwriter
		log.SetOutput(mw)
	} else {
		log.Info("Failed to log to file, using default stderr")
		log.SetOutput(os.Stdout)
	}

	return log
}
