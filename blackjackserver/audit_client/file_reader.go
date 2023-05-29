package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

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

func (auditIO *AuditIO) ReadData() string {
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
