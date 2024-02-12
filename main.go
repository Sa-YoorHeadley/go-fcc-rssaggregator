package main

import (
	// "fmt"

	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("PORT  was not found in the environment")
	}


	// Router Initialization
	router := chi.NewRouter()
	
	// Server Initialization
	server := &http.Server{
		Handler: router,
		Addr: ":" + port,
	}


	// CORS Initialization
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))

	// New Router
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerError)

	router.Mount("/v1", v1Router)


	// Server Start
	log.Printf("Server starting on port: %v", port)
	err := server.ListenAndServe()


	if err != nil {
		log.Fatal(err)
	}
}