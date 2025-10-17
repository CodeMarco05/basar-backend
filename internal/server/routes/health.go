package routes

import (
	"hackathon-basar-backend/internal/config"
	"hackathon-basar-backend/internal/db"
	logger "hackathon-basar-backend/internal/logger"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	databse := config.MongoDB
	err := db.HealthCheck(r.Context(), databse)

	if err != nil {
		log := logger.GetLogger()
		log.Error().Msgf("The DB isn't available %v", err)
		http.Error(w, "The database isn't available.", http.StatusUnauthorized)
	}

	_, err = w.Write([]byte("The application appears to be up and running"))
	if err != nil {
		log := logger.GetLogger()
		log.Error().Msgf("Writing to the host failed during transmitting: %v", err)
	}
}
