package server

import (
	settings "blackjack/config"
	"bufio"
	"encoding/json"
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
}

type PlayIO struct {
}

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

type AuditIO struct {
	filePath string "./audit/d8b952dd-e0cd-4473-81ad-698a858bb21a.log"
}

func (auditIO *AuditIO) ReadData(conn net.Conn) string {
	type Input struct {
		Level string
		Msg   string
		Name  string
		Time  string
	}

	// Open the file
	file, err := os.Open(auditIO.filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	// Create a new scanner
	scanner := bufio.NewScanner(file)

	// Iterate over each line
	for scanner.Scan() {
		// line := scanner.Text()

		fmt.Println(scanner.Bytes())
		// Parse JSON
		var input Input
		err := json.Unmarshal(scanner.Bytes(), &input)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			continue
		}
		fmt.Println(input)

		// Access parsed data
		fmt.Println("Msg:", input.Msg)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
	}
	return "bla"
}

// Do nothing when sending data in audit mode
func (auditIO *AuditIO) SendData(conn net.Conn, msg string) {}

func MakeIO() ioInterface {
	if settings.Mode == settings.PlayMode {
		fmt.Println("Play mode")
		return &PlayIO{}
	} else if settings.Mode == settings.AuditMode {
		fmt.Println("Audit mode")
		return &AuditIO{}
	} else {
		panic("Mode not recognized")
	}
}
