package server

import (
	settings "blackjack/config"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/gobwas/ws/wsutil"
)

type ioInterface interface {
	ReadData(net.Conn) string
	SendData(net.Conn, string)
	GetSeed(net.Conn) int64
}

type PlayIO struct{}

func (playIO *PlayIO) ReadData(conn net.Conn) string {
	// TODO: Improve connection closed vs read timed out error handling
	// Returns empty string if read failed, EOF if connection was closed
	conn.SetReadDeadline(time.Now().Add(settings.ReadTimeout))

	msg, err := wsutil.ReadClientText(conn)
	if err != nil {
		fmt.Println("Read failed")
		if errors.Is(err, io.EOF) {
			fmt.Println("Connection closed by client")
		} else if errors.Is(err, wsutil.ClosedError{Code: 1001}) {
			fmt.Println("Connection closed by client, ws closed")
		} else if errors.Is(err, os.ErrDeadlineExceeded) {
			fmt.Println("Read timed out")
		} else {
			fmt.Printf("Some other error: %e\n", err)
		}
		return "EOF"
	}

	msg_str := string(msg)
	return msg_str
}

func (playIO *PlayIO) SendData(conn net.Conn, msg string) {
	err := wsutil.WriteServerText(conn, []byte(msg))
	if err != nil {
		fmt.Println("Send failed, ", err)
	}
}

func (playIO *PlayIO) GetSeed(net.Conn) int64 {
	if settings.Mode == settings.AuditMode {
		return settings.GetSeed()
	}
	return time.Now().UnixNano()
}

func MakeIO() ioInterface {
	return &PlayIO{}
}
