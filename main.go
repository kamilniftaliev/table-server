package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"

	"github.com/kamilniftaliev/table-server/api"
	"github.com/kamilniftaliev/table-server/api/helpers"
	"github.com/rs/cors"
)

func main() {
	router := chi.NewRouter()

	router.Use(cors.New(cors.Options{
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler)

	router.Use(helpers.Middleware())

	router.Handle("/", handler.Playground("GraphQL playground", "/api"))
	router.Handle("/api", handler.GraphQL(api.NewExecutableSchema(api.Config{
		Resolvers: &api.Resolver{},
	})))

	log.Printf("Serving at http://localhost:3333/")
	log.Fatal(http.ListenAndServe(":3333", router))
}
