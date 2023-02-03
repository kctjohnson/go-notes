package main

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Printf("This app will eventually handle creating and handling all\n")
	fmt.Printf("note related server requests with a sqlite database rather than MySQL.\n")
	fmt.Printf("This will allow the user to just run the server wherever they'd like,\n")
	fmt.Printf("and then they can use the note app like normal, whether the notes are\n")
	fmt.Printf("stored locally or on a server.\n")

	db, err := openDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func openDB() (*sqlx.DB, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	dbPath := "/.config/gonotes/"
	dbName := "gonotes.db"
	dbFullPath := homeDir + dbPath + dbName
	db, err := sqlx.Open("sqlite3", dbFullPath)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
