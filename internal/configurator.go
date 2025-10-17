package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Version string
	Port    uint16
}

var Config *AppConfig

func LoadApplicationConfig() {
	logger := GetLogger()

	err := godotenv.Load()
	if err != nil {
		logger.Error().Msgf("Error loading .env file: %v", err)
		os.Exit(1)
	}

	logger.Info().Msgf("16 bit max value %v")

	Config = &AppConfig{
		Version: os.Getenv("VERSION"),
		Port: func() uint16 {
			res, err := strconv.ParseUint(os.Getenv("PORT"), 10, 16)
			if err != nil {
				logger.Error().Msgf("Error parsing PORT: %v", err)
				os.Exit(1)
			}
			return uint16(res)
		}(),
	}

	// print the application config
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

	//

	return r
}

func Serve(r *chi.Mux) {
	l := GetLogger()

	port := Config.Port

	l.Info().Msg(fmt.Sprintf("Starting on port: %d", port))

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		l.Info().Msgf("Error during application startup %v", err)
	}
}
