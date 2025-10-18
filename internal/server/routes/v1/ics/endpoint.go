package ics

import (
	"encoding/json"
	"hackathon-basar-backend/internal/logger"
	"hackathon-basar-backend/internal/server"
	"net/http"
	"os"
	"sort"
)

func GetAllAvailableFileNames(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	entries, err := os.ReadDir(server.FileStoragePath)
	if err != nil {
		log.Error().Msgf("Something went wrong while reading directory: %s", err)
		http.Error(w, "Something went wrong while reading the storage dir", http.StatusInternalServerError)
		return
	}

	var fileNames []string
	for _, entry := range entries {
		fileNames = append(fileNames, entry.Name())
	}

	sort.Strings(fileNames)
	jsonResponse, err := json.Marshal(fileNames)
	if err != nil {
		log.Error().Msgf("Something went wrong while marshalling response: %s", err)
		http.Error(w, "Something went wrong while marshalling response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(string(jsonResponse))
	if err != nil {
		log.Error().Msgf("Invalid request without creatorId")
		http.Error(w, "Invalid request without creatorId", http.StatusInternalServerError)
		return
	}
}
