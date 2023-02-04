package main

import (
	"fmt"
	"go-notes/cmd/gonotes/models"
	"log"
	"net/http"
	"os"

	"go-notes/internal/config"
	"go-notes/internal/graphql"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/go-sql-driver/mysql"
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
	conf := config.NewConfig()
	conf.Init()

	log.Println(conf.GqlClient.Endpoint)

	// Create the graphql connection client
	client := graphql.NewClient(conf.GqlClient.Endpoint, http.DefaultClient)

	// Start the charm CLI UI
	mainModel := models.NewMain(client)
	program := tea.NewProgram(mainModel)
	if _, err := program.Run(); err != nil {
		fmt.Printf("Failed to run program: %v", err)
		os.Exit(1)
	}
}
