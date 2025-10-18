package posts

import (
	"encoding/json"
	"errors"
	"fmt"
	"hackathon-basar-backend/internal/config"
	"hackathon-basar-backend/internal/db"
	"hackathon-basar-backend/internal/logger"
	"hackathon-basar-backend/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	posts, err := db.GetAllPosts(r.Context(), config.MongoDB)
	if err != nil {
		log.Error().Msgf("Getting all posts failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Set content type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		log.Error().Msgf("Writing to the host failed during transmitting: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func InsertPost(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()
	var post models.InsertPost

	err := json.NewDecoder(r.Body).Decode(&post)

	if err != nil {
		log.Error().Msgf("Invalid request with body: %v", r.Body)
		http.Error(w, "Invalid request body when transforming to the required object.", http.StatusBadRequest)
		return
	}

	// Validate the struct
	validate := validator.New()
	if err := validate.Struct(post); err != nil {
		log.Error().Msgf("Validation failed: %v", err)

		// Get detailed validation errors
		var validationErrors validator.ValidationErrors
		errors.As(err, &validationErrors)
		errorMessages := make(map[string]string)

		for _, fieldError := range validationErrors {
			errorMessages[fieldError.Field()] = fmt.Sprintf("Field '%s' failed validation: %s", fieldError.Field(), fieldError.Tag())
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  "Validation failed",
			"fields": errorMessages,
		})
		if err != nil {
			log.Error().Msgf("Writing to the host failed during transmitting: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		return
	}

	objectIdString, err := db.InsertPost(r.Context(), config.MongoDB, post)
	if err != nil {
		log.Error().Msgf("Inserting post failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resultString := struct {
		ID string `json:"id"`
	}{
		objectIdString,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resultString)
	if err != nil {
		log.Error().Msgf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	postId := chi.URLParam(r, "postId")

	if postId == "" {
		log.Error().Msgf("Invalid request without postId")
		http.Error(w, "PostId was missing or given under a false key", http.StatusBadRequest)
		return
	}

	post, err := db.GetPostByID(r.Context(), config.MongoDB, postId)
	if err != nil {
		log.Error().Msgf("Error during post fetching: %v", err)
		http.Error(w, "PostId was not found or it was an internal server error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(post)

	if err != nil {
		log.Error().Msgf("Writing to the host failed during transmitting: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetAllTags(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	tags, err := db.GetAllTags(r.Context(), config.MongoDB)
	if err != nil {
		log.Error().Msgf("Getting all tags failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(tags)

	if err != nil {
		log.Error().Msgf("Writing to the host failed during transmitting: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
