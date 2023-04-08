package main

import (
	"apiserver/models"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const (
	// host = "localhost"
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
			id VARCHAR(50) PRIMARY KEY,
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
			'f4e483f4-da9c-49aa-aea0-61c416515e39', 'root',
			'$2a$10$9sdxTP5UUq3WyApkEdXrb..H5Vvdzo68IG6eOx0zyj1Xkqk5pVL/W'
		);`,
	)
	if err != nil {
		fmt.Printf("error creating root user %s\n", err)
		return errors.New("Cant create root user")
	}
	fmt.Println("Root user created")
	return nil
}

func (database *UsersDatabase) getUser(username string) (models.User, error) {
	stmt, err := database.db.Prepare("SELECT id, username, password FROM users WHERE username = $1")
	defer stmt.Close()

	// Try using an interface to unpack User struct
	var (
		id          string
		db_username string
		password    string
	)
	err = stmt.QueryRow(username).Scan(&id, &db_username, &password)
	if err != nil {
		fmt.Println(err)
		return models.User{}, errors.New("User not found")
	}

	user := models.User{Id: id, Username: db_username, Password: password}
	return user, nil
}

func (database *UsersDatabase) addUser(username string, password string) (models.User, error) {
	stmt, err := database.db.Prepare(`
		INSERT INTO users (id, username, password) VALUES ($1, $2, $3)`,
	)
	defer stmt.Close()

	id := uuid.NewString()
	hashed_password, _ := HashPassword(password)
	_, err = stmt.Exec(id, username, hashed_password)
	if err != nil {
		return models.User{}, errors.New("Cannot add user")
	}

	user := models.User{Id: id, Username: username, Password: hashed_password}
	return user, nil
}

// func main() {
// 	// var err error
// 	db, err := NewUsersDatabase()
// 	if err != nil {
// 		fmt.Printf("error couldn't connect to database: %s\n", err)
// 	}
// 	defer db.Close()
// 	// db.query()
// 	db.getUserData("root")
// }
