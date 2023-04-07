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

type databaseInterface interface {
	connect()
	query(string)
	Close()
}

type Database struct {
	db *sql.DB
}

func (database *Database) Close() {
	database.db.Close()
}

func (database *Database) query() {
	result, err := database.db.Query("SELECT * from users;")
	if err != nil {
		fmt.Println("query failed")
	}
	defer result.Close()
	fmt.Println(result)
}

func (database *Database) connect() error {
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s sslmode=disable",
		host, port, user, password,
	)

	var (
		tryCount uint
		waitTime time.Duration = initialWait
		err      error
	)
	// Retry when connecting as database container might not be up yet
	for {
		tryCount++
		if tryCount > maxRetries {
			fmt.Println("Connection to database could not be established")
			return errors.New("max retries to database exceeded")
		}

		database.db, err = sql.Open("postgres", psqlconn)

		// check db
		err = database.db.Ping()
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
	return nil
}

type UsersDatabase struct {
	Database
}

func NewUsersDatabase() (UsersDatabase, error) {
	database := UsersDatabase{}
	database.connect()
	if database.checkIfTableExists() {
		return database, nil
	}
	// Table does not exist, create it
	database.createTable()
	database.addRootUser()
	return database, nil
}

func (database *UsersDatabase) checkIfTableExists() bool {
	fmt.Println("Checking if users table exists")
	_, err := database.db.Exec(`SELECT 1 FROM users;`)
	return err == nil
}

func (database *UsersDatabase) createTable() error {
	fmt.Println("Creating users table")
	_, err := database.db.Exec(`
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
	return nil
}

func (database *UsersDatabase) addRootUser() error {
	fmt.Println("Create root user")
	_, err := database.db.Exec(`
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
	fmt.Println("Root user created")
	return nil
}
