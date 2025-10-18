package ics

import (
	"encoding/json"
	"hackathon-basar-backend/internal/logger"
	"hackathon-basar-backend/internal/server"
	"net/http"
	"os"
)

func GetAllAvailableFileNames(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	entries, err := os.ReadDir(server.FileStoragePath)
	if err != nil {
		log.Error().Msgf("Something went wrong while reading directory: %s", err)
		http.Error(w, "Something went wrong while reading the storage dir", http.StatusInternalServerError)
		return
	}

	var responseStruct = struct {
		Filenames []string `json:"filenames"`
	}{}

	for _, entry := range entries {
		responseStruct.Filenames = append(responseStruct.Filenames, entry.Name())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(responseStruct)
	if err != nil {
		log.Error().Msgf("Invalid request without creatorId")
		http.Error(w, "Invalid request without creatorId", http.StatusInternalServerError)
		return
	}
}

func GetIcsFile(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	fileName := r.PathValue("fileName")
	if fileName == "" {
		log.Error().Msgf("Invalid request without fileName")
		http.Error(w, "Invalid request without fileName", http.StatusBadRequest)
		return
	}

	filePath := server.FileStoragePath + "/" + fileName

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Error().Msgf("File %s does not exist", fileName)
		http.Error(w, "File does not exist", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}
