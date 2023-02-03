package main

import (
	"fmt"
	"go-notes/cmd/gonotes/models"
	"log"
	"net/http"
	"os"

	"go-notes/internal/graphql"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
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
	initConfig()
	endpoint := viper.GetString("gql_client.endpoint")

	log.Println(endpoint)

	// Create the graphql connection client
	client := graphql.NewClient(endpoint, http.DefaultClient)

	// Start the charm CLI UI
	mainModel := models.NewMain(client)
	program := tea.NewProgram(mainModel)
	if _, err := program.Run(); err != nil {
		fmt.Printf("Failed to run program: %v", err)
		os.Exit(1)
	}
}

func initConfig() {
	// Set up the config
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.config/gonotes/")
	err := viper.ReadInConfig()

	// If it failed to read in the config, create a new blank one
	if err != nil {
		log.Printf("Creating blank config file at $HOME/.config/gonotes\n")
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		configPath := "/.config/gonotes/"
		configName := "config.json"

		configData, err := os.ReadFile("example_config.json")
		if err != nil {
			panic(err)
		}

		err = os.MkdirAll(homeDir+configPath, os.ModePerm)
		if err != nil {
			panic(err)
		}

		dir := homeDir + configPath + configName
		err = os.WriteFile(dir, configData, 0777)
		if err != nil {
			panic(err)
		}
	}
}
