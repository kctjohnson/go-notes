package main

import (
	"context"
	"fmt"
	// "go-notes/cmd/gonotes/models"
	// "go-notes/internal/db/repositories"
	// "go-notes/internal/services"
	// "log"
	"net/http"
	"os"

	gql "go-notes/cmd/gonotes/graphql"

	"github.com/Khan/genqlient/graphql"
	//tea "github.com/charmbracelet/bubbletea"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func initConfig() {
	// Set up the config
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.config/gonotes/")
	err := viper.ReadInConfig()

	// If it failed to read in the config, create a new blank one
	if err != nil {
		fmt.Printf("Creating blank config file at $HOME/.config/gonotes\n")
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

func main() {
	ctx := context.Background()
	client := graphql.NewClient("http://localhost:3030/graphql", http.DefaultClient)

	_, err := gql.CreateNote(ctx, client, "New Note!")
	if err != nil {
		panic(err)
	}

	res, err := gql.GetNotes(ctx, client)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", res)
	// // Set up logging
	// f, err := tea.LogToFile("debug.log", "debug")
	// if err != nil {
	// 	fmt.Println("fatal: ", err)
	// 	os.Exit(1)
	// }
	// defer f.Close()
	//
	// // Set up env vars
	// initConfig()
	//
	// dbUser := viper.GetString("db.user")
	// dbPassword := viper.GetString("db.password")
	// dbIP := viper.GetString("db.ip")
	// dbPort := viper.GetString("db.port")
	//
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/gonotes?parseTime=true", dbUser, dbPassword, dbIP, dbPort)
	//
	// log.Printf("%s\n", dsn)
	//
	// // Set up the DB
	// db, err := openDB(dsn)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// notesRepo := repositories.NewNotesRepository(db)
	// notesService := services.NewNotesService(notesRepo)
	//
	// // Start the charm CLI UI
	// mainModel := models.NewMain(notesService)
	// program := tea.NewProgram(mainModel)
	// if _, err := program.Run(); err != nil {
	// 	fmt.Printf("Failed to run program: %v", err)
	// 	os.Exit(1)
	// }
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
