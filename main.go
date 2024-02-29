	package main

	import (
		"time"

		"database/sql"
		"log"
		"net/http"
		"os"

		"github.com/Sa-YoorHeadley/go-fcc-rssaggregator/internal/database"
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
		port := os.Getenv("PORT")

		if port == "" {
			log.Fatal("PORT was not found in the environment")
		}

		// Import and Check DB Connection
		dbUrl := os.Getenv("DB_URL")

		if dbUrl == "" {
			log.Fatal("DB_URL was not found in the environment")
		}

		conn, err := sql.Open("postgres", dbUrl)

		if err != nil {
			log.Fatal("Can't connect to databse:", err)
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	go startScraping(db, 10, time.Minute)

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
	
	v1Router.Post("/users", apiCfg.handlerCreateUser)

	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUserByAPIKey))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	router.Mount("/v1", v1Router)

	// Server Start
	log.Printf("Server starting on port: %v", port)
	err = server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}