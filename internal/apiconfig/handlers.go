package apiconfig

import (
	"net/http"
	"log"
	"strconv"
)

func (api *Conf) HandlerApp() http.Handler {
	return http.StripPrefix("/app", http.FileServer(http.Dir(api.FilepathRoot)))
}

func (api *Conf) HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(http.StatusText(http.StatusOK)))
	if err != nil {
		log.Printf("Error writing readiness response: %v", err)
	}
}

func (api *Conf) HandlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hits: " + strconv.Itoa(int(api.FileserverHits.Load()))))
	if err != nil {
		log.Printf("Error writing metrics response: %v", err)
	}
}

func (api *Conf) HandlerReset(w http.ResponseWriter, r *http.Request) {
	api.FileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hits reset to 0"))
	if err != nil {
		log.Printf("Error writing reset response: %v", err)
	}
}