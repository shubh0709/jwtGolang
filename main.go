package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type User struct {
	UserID   string `json:"userid"`
	Password string `json:"password"`
}

var DB *sql.DB

func setupDB() (*sql.DB, error) {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "12345"
		dbname   = "jwtGolang"
	)

	// Create a connection string.
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON data
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the user ID and password into the table
	_, err = DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.UserID, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User ID and password stored successfully")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Return a success response
	w.WriteHeader(http.StatusOK)
}

func main() {
	// Open a connection to the database.
	db, err := setupDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	DB = db

	_, err = db.Exec("SELECT 1 FROM users LIMIT 1")
	if err != nil {
		// Create the users table if it doesn't exist.
		_, err = db.Exec("CREATE TABLE users (id serial PRIMARY KEY, username text NOT NULL, password text NOT NULL)")
		if err != nil {
			panic(err)
		}
	}

	http.HandleFunc("/user", createUserHandler)
	http.HandleFunc("/", healthCheckHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

	fmt.Println("User created successfully!")
}
