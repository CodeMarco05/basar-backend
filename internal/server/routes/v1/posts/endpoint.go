package posts

import (
	"encoding/json"
	"hackathon-basar-backend/internal/config"
	"hackathon-basar-backend/internal/db"
	"hackathon-basar-backend/internal/logger"
	"hackathon-basar-backend/internal/models"
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

func InsertPost(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var post models.InsertPost

	err := json.NewDecoder(r.Body).Decode(&post)

	if err != nil {
		log.Error().Msgf("Invalid request with body: %v", r.Body)
		http.Error(w, "Invalid request body when transforming to the required object.", http.StatusBadRequest)
	}

	objectIdString, err := db.InsertPost(r.Context(), config.MongoDB, post)
	if err != nil {
		log.Error().Msgf("Inserting post failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	resultString := struct {
		ID string `json:"id"`
	}{
		objectIdString,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resultString)
}
func DeletePost(w http.ResponseWriter, r *http.Request) {
	/*log := logger.GetLogger()

	postToDelete := struct {
		ID string `json:"id"`
	}{}*/
}
func UpdatePost(w http.ResponseWriter, r *http.Request) {}
