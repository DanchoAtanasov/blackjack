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

func connectToDatabase() error {
	// connection string
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

	fmt.Println("Connected to database!")
	return nil
}
