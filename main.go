package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/bbarrington0099/Chirpy/internal/apiconfig"
	"github.com/bbarrington0099/Chirpy/internal/database"
	"github.com/bbarrington0099/Chirpy/internal/middleware"
)

func main() {
	godotenv.Load()

	platform := os.Getenv("PLATFORM")

	// Connect to the database
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Load app configuration
	apiConfig := &apiconfig.Conf{
		Port:         "8080",
		FilepathRoot: ".",
		FileserverHits: atomic.Int32{},
		QueryCollection: database.New(db),
		Platform: platform,
	}

	mux := http.NewServeMux()

	// Initialize configuration converted to middleware compatible type
	middlewareApiInstance := (*middleware.LocalConf)(apiConfig)

	// /app/
	mux.Handle("/app/", middlewareApiInstance.MiddlewareFileserverHits(apiConfig.HandlerApp()))

	// /api/
	mux.HandleFunc("GET /api/healthz", apiConfig.HandlerReadiness)
	mux.HandleFunc("POST /api/users", apiConfig.HandlerCreateUser)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.HandlerGetChirpByID)
	mux.HandleFunc("GET /api/chirps", apiConfig.HandlerGetChirps)
	mux.HandleFunc("POST /api/chirps", apiConfig.HandlerCreateChirp)
	
	// /admin/
	mux.HandleFunc("GET /admin/metrics", apiConfig.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiConfig.HandlerReset)

	srv := &http.Server{
		Addr:    ":" + apiConfig.Port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", apiConfig.FilepathRoot, apiConfig.Port)
	log.Fatal(srv.ListenAndServe())
}