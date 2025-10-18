package server

import (
	"fmt"
	"hackathon-basar-backend/internal/config"
	"hackathon-basar-backend/internal/logger"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var availableFiles = []string{
	"A23a_5", "A24a_3", "A23b_5",
	"B23a_5",
	"I23a_5", "I23c_5",
	"T23a_5", "T24a_3",
	"W23a_5", "W23c_5",
	"W24a_3", "W24b_3", "W24c_3", "W24d_3", "W24e_3",
}

var FileStoragePath = "./ics-files"

func IcsFetcherSetupAndStart() {
	log := logger.GetLogger()

	// remove the storage path if it exists
	err := os.RemoveAll(FileStoragePath)
	if err != nil {
		log.Error().Msgf("Error removing files: %v", err)
		return
	}

	// create a local path storage path
	err = os.MkdirAll(FileStoragePath, os.ModePerm)
	if err != nil {
		log.Error().Msgf("Error creating directory: %v", err)
		os.Exit(1)
	}

	// Start ticker in a goroutine
	go func() {
		i := fmt.Sprintf("%vs", config.Config.IcsFileScrapeIntervalInSeconds)
		duration, err := time.ParseDuration(i)

		if err != nil {
			log.Error().Msgf("Error parsing IcsFileScrapeIntervalInSeconds: %v", err)
			os.Exit(1)
		}

		ticker := time.NewTicker(duration)
		defer ticker.Stop()

		// run it one time at startup
		icsFetcher()

		for {
			select {
			case <-ticker.C:
				icsFetcher()
			}
		}
	}()
}

func icsFetcher() {
	log := logger.GetLogger()
	for _, entry := range availableFiles {
		url := fmt.Sprintf("https://cis.nordakademie.de/fileadmin/Infos/Stundenplaene/%v.ics", entry)

		resp, err := http.Get(url)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Error().Msgf("Error closing body: %v", err)
				return
			}
		}(resp.Body)
		if err != nil {
			log.Error().Msgf("Error fetching %v: %v", entry, err)
			continue
		}

		// Check if the request was successful
		if resp.StatusCode != http.StatusOK {
			log.Error().Msgf("Error fetching %v: %v", entry, resp.StatusCode)
			return
		}

		// Full path to save file
		entryWithFileExtension := entry + ".ics"
		fullPath := filepath.Join(FileStoragePath, entryWithFileExtension)

		outFile, err := os.Create(fullPath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}

		_, err = io.Copy(outFile, resp.Body)
		if err != nil {
			log.Error().Msgf("Error reading %v: %v", entry, err)
			continue
		}

		err = outFile.Close()
		if err != nil {
			log.Error().Msgf("Error closing file: %v", err)
			return
		}

		log.Info().Msgf("Successfully fetched %v", entry)
	}
}
