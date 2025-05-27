package apiconfig

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(fmt.Appendf(nil, 
		`<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>`, 
		api.FileserverHits.Load(),
	))
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

func (api *Conf) HandlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Chirp string `json:"body"`
	}

	profaneWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := decoder.Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "Invalid request body"}`))
		if err != nil {
			log.Printf("Error writing error response: %v", err)
		}
		return
	}

	if len(params.Chirp) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "Chirp is too long"}`))
		if err != nil {
			log.Printf("Error writing error response: %v", err)
		}
		return
	}

	cleanedChirp := params.Chirp
	words := strings.Split(cleanedChirp, " ")
	for i, word := range words {
		for _, profane := range profaneWords {
			if strings.EqualFold(word, profane) {
				words[i] = "****"
			}		
		}
	}

	cleanedChirp = strings.Join(words, " ")

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"cleaned_body": "` + cleanedChirp + `"}`))
	if err != nil {
		log.Printf("Error writing success response: %v", err)
	}
}