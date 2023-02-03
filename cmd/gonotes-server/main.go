package main

import (
	gographql "go-notes/cmd/gonotes-server/graphql"
	"go-notes/pkg/db/repositories"
	"go-notes/pkg/services"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/introspection"
)

type Application struct {
	gqlServer *gographql.GQLServer
}

func main() {
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

	http.Handle("/graphql", corsHandler(graphql.HTTPHandler(schema)))
	http.ListenAndServe(":3030", nil)
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
