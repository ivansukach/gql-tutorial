package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-pg/pg/v10"
	"github.com/ivansukach/gql-tutorial/domain"
	customMiddleware "github.com/ivansukach/gql-tutorial/middleware"
	"github.com/ivansukach/gql-tutorial/repository"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ivansukach/gql-tutorial/graph"
	"github.com/ivansukach/gql-tutorial/graph/generated"
	"github.com/joho/godotenv"
)

const defaultPort = "8081"

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	db := repository.New(&pg.Options{User: "su", Password: "su", Database: "gql-tutorial"})
	usersRepo := repository.UsersRepo{DB: db}
	meetupsRepo := repository.MeetupsRepo{DB: db}
	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(customMiddleware.AuthMiddleware(usersRepo))
	defer db.Close()
	db.AddQueryHook(repository.DBLogger{})

	d := domain.NewDomain(usersRepo, meetupsRepo)
	c := generated.Config{Resolvers: &graph.Resolver{
		Domain: d,
	}}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", graph.DataLoaderMiddleware(db, srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
