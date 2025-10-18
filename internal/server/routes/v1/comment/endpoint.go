package comment

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

func CreateComment(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	postId := chi.URLParam(r, "postId")

	if postId == "" {
		log.Error().Msgf("Invalid request without postId")
		http.Error(w, "postId is required", http.StatusBadRequest)
		return
	}

	var comment models.Comment

	err := json.NewDecoder(r.Body).Decode(&comment)

	if err != nil {
		log.Error().Msgf("Invalid request without comment")
		http.Error(w, "comment is required", http.StatusBadRequest)
		return
	}

	// Validate the struct
	validate := validator.New()
	if err := validate.Struct(comment); err != nil {
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

	err = db.AddCommentToPost(r.Context(), config.MongoDB, postId, comment)
	if err != nil {
		log.Error().Msgf("Failed to create comment: %v", err)
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}
}

func DeleteCommentByPostIdCommentCreatorIdCommentId(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	postId := chi.URLParam(r, "postId")
	if postId == "" {
		log.Error().Msgf("Invalid request without postId")
		http.Error(w, "postId is required", http.StatusBadRequest)
		return
	}

	commentCreatorId := chi.URLParam(r, "commentCreatorId")
	if commentCreatorId == "" {
		log.Error().Msgf("Invalid request without commentCreatorId")
		http.Error(w, "commentCreatorId is required", http.StatusBadRequest)
		return
	}

	commentId := chi.URLParam(r, "commentId")
	if commentId == "" {
		log.Error().Msgf("Invalid request without commentId")
		http.Error(w, "commentId is required", http.StatusBadRequest)
		return
	}

}
