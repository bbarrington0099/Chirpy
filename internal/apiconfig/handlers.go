package apiconfig

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/bbarrington0099/Chirpy/internal/database"
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
	if api.Platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		if _, err := w.Write([]byte(`{"error": "Reset is only allowed in development mode"}`)); err != nil {
			log.Printf("Error writing error response: %v", err)
		}
		return
	}
	
	api.FileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte("Hits reset to 0")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(`{"error": "Failed to reset hits"}`)); err != nil {
			log.Printf("Error writing error response: %v", err)
		}
		return
	}
	
	if err := api.QueryCollection.DeleteAllUsers(r.Context()); err != nil {
		log.Printf("Error deleting all users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(`{"error": "Failed to delete all users"}`)); err != nil {
			log.Printf("Error writing error response: %v", err)
		}
		return
	}
	
	w.WriteHeader(http.StatusOK)
}

func (api *Conf) HandlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Chirp string `json:"body"`
		User uuid.UUID `json:"user_id"`
	}

	type ChirpResponse struct {
		ID        uuid.UUID `json:"id"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
		CreatedAt string    `json:"created_at"`
		UpdatedAt string    `json:"updated_at"`
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

	chirpParams := database.CreateChirpParams{
		UserID: params.User,
		Body:   cleanedChirp,
	}

	if chirp, err := api.QueryCollection.CreateChirp(r.Context(), chirpParams); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(`{"error": "Failed to create chirp"}`)); err != nil {
			log.Printf("Error writing error response: %v", err)
		}
		return
	} else {
		chirpResponse := ChirpResponse{
			ID:        chirp.ID,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
			CreatedAt: chirp.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: chirp.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		response, err := json.Marshal(chirpResponse)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error marshalling chirp response: %v", err)
			return
		}
		w.WriteHeader(http.StatusCreated)	
		if _, err := w.Write(response); err != nil {
			log.Printf("Error writing chirp response: %v", err)
			return
		}
	}
}

func (api *Conf) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type userRequest struct {
		Email string `json:"email"`
	}

	type UserResponse struct {
		ID uuid.UUID `json:"id"`
		Email string `json:"email"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	userArgs := userRequest{}
	if err := json.NewDecoder(r.Body).Decode(&userArgs); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "Invalid request body"}`))
		if err != nil {
			log.Printf("Error writing error response: %v", err)
		}
		return
	}

	if userArgs.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "Email is required"}`))
		if err != nil {
			log.Printf("Error writing error response: %v", err)
		}
		return
	}

	user, err := api.QueryCollection.CreateUser(r.Context(), userArgs.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error creating user: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	userResponse := UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	response, _ := json.Marshal(userResponse)
	w.Write(response)
}

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

func (api *Conf) HandlerGetChirps(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	chirps, err := api.QueryCollection.GetChirps(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(`{"error": "Failed to retrieve chirps"}`)); err != nil {
			log.Printf("Error writing error response: %v", err)
		}
		return
	}

	chirpResponse := make([]ChirpResponse, len(chirps))
	for i, chirp := range chirps {
		chirpResponse[i] = ChirpResponse{
			ID:        chirp.ID,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
			CreatedAt: chirp.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: chirp.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(chirpResponse)
	w.Write(response)
}

func (api *Conf) HandlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(`{"error": "Chirp ID is required"}`)); err != nil {
			log.Printf("Error writing error response: %v", err)
		}
		return
	}

	chirp, err := api.QueryCollection.GetChirpById(r.Context(), uuid.MustParse(chirpID))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte(`{"error": "Chirp not found"}`)); err != nil {
			log.Printf("Error writing error response: %v", err)
		}
		return
	}

	chirpResponse := ChirpResponse{
		ID:        chirp.ID,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
		CreatedAt: chirp.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: chirp.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(chirpResponse)
	w.Write(response)
}