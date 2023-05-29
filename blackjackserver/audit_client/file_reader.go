package main

import (
	"blackjack/server"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type FileIO struct {
	file    *os.File
	scanner *bufio.Scanner
}

func MakeFileIO(filePath string) *FileIO {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		panic("File cannot be opened")
	}

	// Create a new scanner
	scanner := bufio.NewScanner(file)

	auditIO := FileIO{file: file, scanner: scanner}
	return &auditIO
}

func (fileIO *FileIO) Close() {
	fileIO.file.Close()
}

type SessionIO struct{ *FileIO }

func MakeSessIO(filePath string) *SessionIO {
	return &SessionIO{MakeFileIO(filePath)}
}

type ConnIO struct{ *FileIO }

func MakeConnIO(filePath string) *ConnIO {
	return &ConnIO{MakeFileIO(filePath)}
}

func (connIO *ConnIO) ReadData() string {
	type Input struct {
		Level string
		Msg   string
		Name  string
		Time  string
	}

	if !connIO.scanner.Scan() {
		if err := connIO.scanner.Err(); err != nil {
			fmt.Println("Error scanning file:", err)
			panic("Scan failed")
		}
	}

	var input Input
	err := json.Unmarshal(connIO.scanner.Bytes(), &input)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		panic("Json not parsed")
	}
	fmt.Println(input)

	// Access parsed data
	fmt.Println("Msg:", input.Msg)

	return input.Msg
}

type SessInput struct {
	Level         string
	Msg           string
	Time          string
	Action        string
	PlayerDetails server.PlayerDetails
	SessionId     string
	Type          string
	Seed          int
	Round         int
}

func (sessIO *SessionIO) ReadData() (SessInput, error) {
	var err error = nil
	if !sessIO.scanner.Scan() {
		if err = sessIO.scanner.Err(); err != nil {
			fmt.Println("Error scanning file:", err)
			return SessInput{}, err
		}
		return SessInput{}, io.EOF
	}

	var input SessInput
	err = json.Unmarshal(sessIO.scanner.Bytes(), &input)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		fmt.Println(sessIO.scanner.Bytes())
		panic("Json not parsed")
	}
	fmt.Println(input)

	return input, err
}
