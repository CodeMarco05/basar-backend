package users

import (
	"encoding/json"
	"hackathon-basar-backend/internal/config"
	"hackathon-basar-backend/internal/db"
	"hackathon-basar-backend/internal/logger"
	"hackathon-basar-backend/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetPostsByCreator(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	creatorId := chi.URLParam(r, "creatorId")

	if creatorId == "" {
		log.Error().Msgf("Invalid request without creatorId")
		http.Error(w, "creatorId was missing or given under a false key", http.StatusBadRequest)
		return
	}

	posts, err := db.GetPostsByCreator(r.Context(), config.MongoDB, creatorId)
	if err != nil {
		log.Error().Msgf("Error during post fetching: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		log.Error().Msgf("Writing to the host failed during transmitting: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func UpdatePostByIdAndCreator(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	postId := chi.URLParam(r, "postId")
	if postId == "" {
		log.Error().Msgf("Invalid request without postId")
		http.Error(w, "postId was missing or given under a false key", http.StatusBadRequest)
		return
	}
	creatorId := chi.URLParam(r, "creatorId")
	if creatorId == "" {
		log.Error().Msgf("Invalid request without creatorId")
		http.Error(w, "creatorId was missing or given under a false key", http.StatusBadRequest)
		return
	}

	var post models.InsertPost

	err := json.NewDecoder(r.Body).Decode(&post)

	if err != nil {
		log.Error().Msgf("Invalid request with body: %v", r.Body)
		http.Error(w, "Invalid request body when transforming to the required object.", http.StatusBadRequest)
		return
	}

}
