package posts

import (
	"encoding/json"
	"hackathon-basar-backend/internal/config"
	"hackathon-basar-backend/internal/db"
	"hackathon-basar-backend/internal/logger"
	"net/http"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	posts, err := db.GetAllPosts(r.Context(), config.MongoDB)
	if err != nil {
		log.Error().Msgf("Getting all posts failed: %v", err)
		http.Error(w, "Unauthorized. The user has no access to Firebase", http.StatusInternalServerError)
	}
	// Set content type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		log.Error().Msgf("Writing to the host failed during transmitting: %v", err)
		http.Error(w, "Unauthorized. The user has no access to Firebase", http.StatusInternalServerError)
	}
}

func InsertPost(w http.ResponseWriter, r *http.Request) {}
func DeletePost(w http.ResponseWriter, r *http.Request) {}
func UpdatePost(w http.ResponseWriter, r *http.Request) {}
