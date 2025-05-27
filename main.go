package main

import (
	"log"
	"net/http"
	"sync/atomic"
	
	"github.com/bbarrington0099/Chirpy/internal/apiconfig"
	"github.com/bbarrington0099/Chirpy/internal/middleware"
)

func main() {
	apiConfig := &apiconfig.Conf{
		Port:         "8080",
		FilepathRoot: ".",
		FileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()

	middlewareApiInstance := (*middleware.LocalConf)(apiConfig)

	// /app/
	mux.Handle("/app/", middlewareApiInstance.MiddlewareFileserverHits(apiConfig.HandlerApp()))

	// /api/
	mux.HandleFunc("GET /api/healthz", apiConfig.HandlerReadiness)
	
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