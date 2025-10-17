package internal

import (
	"encoding/json"
	"fmt"
	"hackathon-basar-backend/internal/db"
	"hackathon-basar-backend/internal/logger"
	"hackathon-basar-backend/internal/server/routes"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppConfig struct {
	Version string
	Port    uint16
}

var Config *AppConfig
var MongoDB *mongo.Database

func LoadApplicationConfig() {
	// ------------------------------------------------------------ //
	// load the env config
	logger := logger.GetLogger()

	err := godotenv.Load()
	if err != nil {
		logger.Error().Msgf("Error loading .env file: %v", err)
		os.Exit(1)
	}

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

	// ------------------------------------------------------------ //
	// load the firebase config
	logger.Info().Msg("Start connection to firebase")
	err = InitFirebase()

	if err != nil {
		logger.Error().Msgf("Error initializing firebase: %v", err)
		os.Exit(1)
	}

	logger.Info().Msg("Connection to firebase finished")

	// ------------------------------------------------------------ //
	// load the mongodb connection

	MongoDB, err = db.InitMongoDB("mongodb://localhost:27017", "bazzar")
	if err != nil {
		logger.Error().Msgf("Error setting up MongoDB connection: %v", err)
		os.Exit(1)
	}

}

func ChiConfig() *chi.Mux {
	r := chi.NewRouter()

	// TODO implement own logger middle ware

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(AuthMiddleware)
	r.Get("/health", routes.HealthCheck)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(AuthMiddleware)

	})

	return r
}

func Serve(r *chi.Mux) {
	l := logger.GetLogger()

	port := Config.Port

	l.Info().Msg(fmt.Sprintf("Starting on port: %d", port))

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		l.Info().Msgf("Error during application startup %v", err)
	}
}
