package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
	table "github.com/kamilniftaliev/table-server/graphql"
)

func main() {
	router := chi.NewRouter()

	router.Use(table.Middleware())

	router.Handle("/", handler.Playground("GraphQL playground", "/api"))
	router.Handle("/api", handler.GraphQL(table.NewExecutableSchema(table.Config{Resolvers: &table.Resolver{}})))

	log.Printf("Serving at http://localhost:4444/")
	log.Fatal(http.ListenAndServe(":4444", router))
}
