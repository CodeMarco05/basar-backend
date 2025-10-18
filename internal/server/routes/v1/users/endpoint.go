package users

import (
	"encoding/json"
	"fmt"
	"hackathon-basar-backend/internal/config"
	"hackathon-basar-backend/internal/db"
	"hackathon-basar-backend/internal/logger"
	"hackathon-basar-backend/internal/models"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
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

func PatchPostByIdAndCreatorId(w http.ResponseWriter, r *http.Request) {
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

	// Validate the struct
	validate := validator.New()
	if err := validate.Struct(post); err != nil {
		log.Error().Msgf("Validation failed: %v", err)

		// Get detailed validation errors
		validationErrors := err.(validator.ValidationErrors)
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

	err = db.PatchPostByIdAndCreatorId(r.Context(), config.MongoDB, postId, creatorId, post)
	if err != nil {
		log.Error().Msgf("Error during post update: %v", err)
		errString := fmt.Sprintf("Error during post update: %v", err)
		http.Error(w, errString, http.StatusInternalServerError)
		return
	}
}

func DeletePostByIdAndCreatorId(w http.ResponseWriter, r *http.Request) {
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

	err := db.DeletePostByIdAndCreatorId(r.Context(), config.MongoDB, postId, creatorId)
	if err != nil {
		log.Error().Msgf("Error during post delete: %v", err)
		errString := fmt.Sprintf("Error during post delete: %v", err)
		http.Error(w, errString, http.StatusInternalServerError)
		return
	}
}
