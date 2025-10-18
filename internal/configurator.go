package internal

import (
	"encoding/json"
	"fmt"
	"hackathon-basar-backend/internal/config"
	"hackathon-basar-backend/internal/db"
	"hackathon-basar-backend/internal/logger"
	"hackathon-basar-backend/internal/server/routes"
	"hackathon-basar-backend/internal/server/routes/v1/posts"
	"hackathon-basar-backend/internal/server/routes/v1/users"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func LoadApplicationConfig() {
	// ------------------------------------------------------------ //
	// load the env config
	log := logger.GetLogger()

	err := godotenv.Load()
	if err != nil {
		log.Error().Msgf("Error loading .env file: %v", err)
		os.Exit(1)
	}

	config.Config = &config.AppConfig{
		Version: func() string {
			res := os.Getenv("VERSION")
			if res == "" {
				log.Error().Msg("VERSION environment variable is not set")
				os.Exit(1)
			}
			return res
		}(),
		Port: func() uint16 {
			res, err := strconv.ParseUint(os.Getenv("PORT"), 10, 16)
			if err != nil {
				log.Error().Msgf("Error parsing PORT: %v", err)
				os.Exit(1)
			}
			return uint16(res)
		}(),
		MongoDbURI: func() string {
			res := os.Getenv("MONGODB_URI")
			if res == "" {
				log.Error().Msg("MONGODB_URI environment variable is not set")
				os.Exit(1)
			}
			return res
		}(),
	}

	// print the application config
	configJson, err := json.Marshal(*config.Config)
	if err != nil {
		log.Error().Msgf("Error marshalling config: %v", err)
		os.Exit(1)
	}
	log.Info().Msgf("Config: %v", string(configJson))

	// ------------------------------------------------------------ //
	// load the firebase config
	log.Info().Msg("Start connection to firebase")
	err = InitFirebase()

	if err != nil {
		log.Error().Msgf("Error initializing firebase: %v", err)
		os.Exit(1)
	}

	log.Info().Msg("Connection to firebase finished")

	// ------------------------------------------------------------ //
	// load the mongodb connection

	config.MongoDB, err = db.InitMongoDB(config.Config.MongoDbURI, "bazzar")
	if err != nil {
		log.Error().Msgf("Error setting up MongoDB connection: %v", err)
		os.Exit(1)
	}

}

func ChiConfig() *chi.Mux {
	r := chi.NewRouter()

	// TODO implement own logger middle ware

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	//r.Use(AuthMiddleware)
	r.Get("/health", routes.HealthCheck)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/posts", func(r chi.Router) {
			r.Get("/", posts.GetAllPosts)
			r.Get("/{postId}", posts.GetPost)
			r.Post("/", posts.InsertPost)
		})

		r.Route("/users", func(r chi.Router) {
			r.Get("/{creatorId}/posts", users.GetPostsByCreator)
			r.Patch("/{creatorId}/{postId}", users.PatchPostByIdAndCreatorId)
			r.Delete("/{creatorId}/{postId}", users.DeletePostByIdAndCreatorId)
		})
	})

	return r
}

func Serve(r *chi.Mux) {
	l := logger.GetLogger()

	port := config.Config.Port

	l.Info().Msg(fmt.Sprintf("Starting on port: %d", port))

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		l.Info().Msgf("Error during application startup %v", err)
	}
}
