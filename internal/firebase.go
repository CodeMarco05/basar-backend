package internal

import (
	"context"
	"hackathon-basar-backend/internal/logger"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var authClient *auth.Client

func InitFirebase() error {
	logger := logger.GetLogger()

	// get current dir
	cwd, err := os.Getwd()
	if err != nil {
		logger.Fatal().Err(err).Msg("could not get current working directory")
		os.Exit(1)
	}

	configPath := filepath.Join(cwd, "serviceAccountKey.json")

	ctx := context.Background()
	opt := option.WithCredentialsFile(configPath)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return err
	}

	authClient, err = app.Auth(ctx)
	return err
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logger.GetLogger()

		// Extract token from header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Info().Msg("Authorization header missing")
			http.Error(w, "Unauthorized. The auth header is missing", http.StatusUnauthorized)
			return
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify token
		token, err := verifyToken(idToken)
		if err != nil {
			logger.Info().Err(err).Msg("Unauthorized access to Firebase API endpoint detected")
			http.Error(w, "Unauthorized. The user has no access to Firebase", http.StatusUnauthorized)
			return
		}

		// Token is valid - add user info to context and continue
		ctx := context.WithValue(r.Context(), "userID", token.UID)
		ctx = context.WithValue(r.Context(), "userToken", token)

		// Pass to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func verifyToken(idToken string) (*auth.Token, error) {
	ctx := context.Background()
	return authClient.VerifyIDToken(ctx, idToken)
}
