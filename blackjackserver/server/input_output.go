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
	"strconv"
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
	return 1685355255139882952
	// return time.Now().UnixNano()
}

type AuditIO struct {
	file    *os.File
	scanner *bufio.Scanner
}

func MakeAuditIO(filePath string) *AuditIO {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		panic("File cannot be opened")
	}

	// Create a new scanner
	scanner := bufio.NewScanner(file)

	auditIO := AuditIO{file: file, scanner: scanner}
	return &auditIO
}

func (auditIO *AuditIO) Close() {
	auditIO.file.Close()
}

func (auditIO *AuditIO) ReadData(conn net.Conn) string {
	type Input struct {
		Level string
		Msg   string
		Name  string
		Time  string
	}

	if !auditIO.scanner.Scan() {
		if err := auditIO.scanner.Err(); err != nil {
			fmt.Println("Error scanning file:", err)
			panic("Scan failed")
		}
	}

	var input Input
	err := json.Unmarshal(auditIO.scanner.Bytes(), &input)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		panic("Json not parsed")
	}
	fmt.Println(input)

	// Access parsed data
	fmt.Println("Msg:", input.Msg)

	return input.Msg
}

// Do nothing when sending data in audit mode
func (auditIO *AuditIO) SendData(conn net.Conn, msg string) {
	err := wsutil.WriteServerText(conn, []byte(msg))
	if err != nil {
		fmt.Println("Send failed, ", err)
	}
}

func (auditIO *AuditIO) GetSeed(conn net.Conn) int64 {
	msg := auditIO.ReadData(conn)
	seed, err := strconv.ParseInt(msg, 10, 64)
	if err != nil {
		panic("cant read seed")
	}
	return seed
}

func MakeIO() ioInterface {
	if settings.Mode == settings.PlayMode {
		fmt.Println("Play mode")
		return &PlayIO{}
	} else if settings.Mode == settings.AuditMode {
		fmt.Println("Audit mode")
		return MakeAuditIO(settings.AuditLogFile)
	} else {
		panic("Mode not recognized")
	}
}
