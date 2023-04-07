package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "database"
	port     = 5432
	user     = "postgres"
	password = "superpassword"
	// dbname   = "MyDatabase"

	// TODO: to be tuned
	maxRetries      = 5
	initialWait     = 200 * time.Millisecond
	maxWaitInterval = 3 * time.Second
)

// TODO improve database abstraction
func connectToDatabase() error {
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s sslmode=disable",
		host, port, user, password,
	)

	var (
		tryCount uint
		waitTime time.Duration = initialWait
		db       *sql.DB
		err      error
	)
	for {
		tryCount++
		if tryCount > maxRetries {
			fmt.Println("Connection to database could not be established")
			return errors.New("max retries to database exceeded")
		}

		db, err = sql.Open("postgres", psqlconn)

		// check db
		err = db.Ping()
		if err == nil {
			break
		}

		fmt.Println("retrying")
		// Backoff
		time.Sleep(waitTime)
		waitTime *= 2
		if waitTime > maxWaitInterval {
			waitTime = maxWaitInterval
		}
	}

	// close database
	defer db.Close()
	fmt.Println("Connected")

	fmt.Println("Checking if users table exists")
	_, err = db.Exec(`SELECT 1 FROM users;`)
	if err == nil {
		fmt.Println("Users table exists, conneciton complete")
		return nil
	}

	fmt.Println("Users table does not exist, creating...")
	_, err = db.Exec(`
		CREATE TABLE users(
			id INT PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(100) NOT NULL
		);`,
	)
	if err != nil {
		fmt.Printf("error creating table %s\n", err)
		return errors.New("Cant create table")
	}
	fmt.Println("Users table created")

	fmt.Println("Create root user")
	_, err = db.Exec(`
		INSERT INTO users(
			id, username, password
		) VALUES (
			1, 'root', 'pass'
		);`,
	)
	if err != nil {
		fmt.Printf("error creating root user %s\n", err)
		return errors.New("Cant create root user")
	}

	fmt.Println("Connected to database!")
	return nil
}
