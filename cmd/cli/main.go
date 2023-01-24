package main

import (
	"fmt"
	"go-notes/cmd/cli/models"
	"go-notes/pkg/db/repositories"
	"go-notes/pkg/services"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	// Set up logging
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal: ", err)
		os.Exit(1)
	}
	defer f.Close()

	// Set up env vars
	err = godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	dsn := os.Getenv("DSN")
	log.Printf("%s\n", dsn)

	// Set up the DB
	db, err := openDB(dsn)
	if err != nil {
		panic(err)
	}

	notesRepo := repositories.NewNotesRepository(db)
	notesService := services.NewNotesService(notesRepo)

	// Start the charm CLI UI
	mainModel := models.NewMain(notesService)
	program := tea.NewProgram(mainModel) //, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Failed to run program: %v", err)
		os.Exit(1)
	}
}

func openDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
