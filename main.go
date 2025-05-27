package main

import (
	"log"
	"net/http"
	"sync/atomic"
	
	"github.com/bbarrington0099/Chirpy/internal/apiConfig"
)

func main() {
	apiConfig := &apiConfig.Conf{
		Port:         "8080",
		FilepathRoot: ".",
		FileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiConfig.middlewareFileserverHits(apiConfig.HandlerApp()))
	mux.HandleFunc("GET /healthz", apiConfig.HandlerReadiness)
	mux.HandleFunc("GET /metrics", apiConfig.HandlerMetrics)
	mux.HandleFunc("POST /reset", apiConfig.HandlerReset)

	srv := &http.Server{
		Addr:    ":" + apiConfig.Port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", apiConfig.FilepathRoot, apiConfig.Port)
	log.Fatal(srv.ListenAndServe())
}