package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/digitalhumanlife/rss-aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in ENV.")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in ENV.")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cannot connect to DB: ", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Lnik"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1router := chi.NewRouter()
	v1router.Get("/readiness", handlerReadiness)
	v1router.Get("/err", handlerError)
	v1router.Post("/users", apiCfg.handlerCreateUser)

	router.Mount("/v1", v1router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server is running on port %s", portString)
	err1 := srv.ListenAndServe()
	if err1 != nil {
		log.Fatal(err1)
	}

	fmt.Println("PORT is ", portString)
}
