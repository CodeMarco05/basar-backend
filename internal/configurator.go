package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ocm-go/internal/logging"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Version string
	Port    int16
}

var Config *AppConfig

func LoadApplicationConfig() {
	logger := logging.GetLogger()

	err := godotenv.Load()
	if err != nil {
		logger.Error().Msgf("Error loading .env file: %v", err)
		os.Exit(1)
	}

	Config = &AppConfig{
		Version: os.Getenv("VERSION"),
		Port: func() int16 {
			res, err := strconv.ParseInt(os.Getenv("PORT"), 10, 16)
			if err != nil {
				logger.Error().Msgf("Error parsing PORT: %v", err)
				os.Exit(1)
			}
			return int16(res)
		}(),
	}

	// prin the application config
	configJson, err := json.Marshal(*Config)
	if err != nil {
		logger.Error().Msgf("Error marshalling config: %v", err)
		os.Exit(1)
	}
	logger.Info().Msgf("Config: %v", string(configJson))
}

func ChiConfig() *chi.Mux {
	r := chi.NewRouter()

	// TODO implement own logger middle ware
	// r.Use(middleware.Logger)

	RegisterAllEndpoints(r)

	return r
}

func Serve(r *chi.Mux) {
	l := logging.GetLogger()

	port := 3000

	l.Info().Msg(fmt.Sprintf("Starting on port: %d", port))

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		l.Info().Msgf("Error during application startup %v", err)
	}
}
