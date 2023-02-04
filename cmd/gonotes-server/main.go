package main

import (
	"fmt"
	gographql "go-notes/cmd/gonotes-server/graphql"
	"go-notes/internal/config"
	"go-notes/internal/db/repositories"
	"go-notes/internal/services"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/introspection"
)

func main() {
	// Init the config
	conf := config.NewConfig()
	conf.Init()

	db, err := openDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	nr := repositories.NewNotesRepository(db)
	ns := services.NewNotesService(nr)
	gql := gographql.GQLServer{
		NotesGql: &gographql.NotesGql{
			NotesService: ns,
		},
	}

	schema := gql.Schema()
	introspection.AddIntrospectionToSchema(schema)

	fmt.Printf("Starting graphql server on port %s\n", conf.GqlServer.Port)
	http.Handle("/graphql", corsHandler(graphql.HTTPHandler(schema)))
	http.ListenAndServe(":"+conf.GqlServer.Port, nil)
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

func corsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		h.ServeHTTP(w, r)
	}
}
