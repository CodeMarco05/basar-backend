package db

import (
	"context"
	"fmt"
	"hackathon-basar-backend/internal/logger"
	"hackathon-basar-backend/internal/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitMongoDB initializes and returns a MongoDB client with connection pooling and retry logic
func InitMongoDB(uri string, databaseName string) (*mongo.Database, error) {
	logger := logger.GetLogger()

	// Set client options with connection pooling
	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(50).                   // Maximum number of connections in pool
		SetMinPoolSize(5).                    // Minimum number of connections in pool
		SetMaxConnIdleTime(30 * time.Minute). // Close connections after 30 minutes of inactivity
		SetServerSelectionTimeout(10 * time.Second).
		SetSocketTimeout(30 * time.Second).
		SetConnectTimeout(10 * time.Second)

	// Connect to MongoDB with retry logic
	maxRetries := 10
	retryDelay := 2 * time.Second

	var client *mongo.Client
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		logger.Info().Msgf("Attempting to connect to MongoDB (attempt %d/%d)...", attempt, maxRetries)

		client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			logger.Error().Msgf("Failed to connect to MongoDB (attempt %d): %v", attempt, err)
			if attempt < maxRetries {
				time.Sleep(retryDelay)
				continue
			}
			return nil, fmt.Errorf("failed to connect to MongoDB after %d attempts: %w", maxRetries, err)

		}

		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Printf("Failed to ping MongoDB (attempt %d): %v", attempt, err)
			if attempt < maxRetries {
				time.Sleep(retryDelay)
				continue
			}
			return nil, fmt.Errorf("failed to ping MongoDB after %d attempts: %w", maxRetries, err)
		}

		log.Println("Successfully connected to MongoDB!")
		break
	}

	// Get database handle
	database := client.Database(databaseName)

	// Create collections if they don't exist
	if err := ensureCollectionsExist(context.TODO(), database); err != nil {
		return nil, fmt.Errorf("failed to ensure collections exist: %w", err)
	}

	return database, nil
}

// HealthCheck performs a health check on the MongoDB database
func HealthCheck(ctx context.Context, db *mongo.Database) error {
	if db == nil {
		return fmt.Errorf("MongoDB database is nil")
	}

	// Ping the MongoDB server
	if err := db.Client().Ping(ctx, nil); err != nil {
		return fmt.Errorf("MongoDB health check failed: %w", err)
	}
	return nil
}

// ensureCollectionsExist creates required collections if they don't already exist
func ensureCollectionsExist(ctx context.Context, db *mongo.Database) error {
	// List existing collections
	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	log.Printf("Existing collections: %v", collections)

	// Create a map for quick lookup
	existingCollections := make(map[string]bool)
	for _, name := range collections {
		existingCollections[name] = true
	}

	// Define required collections
	requiredCollections := []string{
		"posts",
		// Add more collection names as needed
	}

	// Create missing collections
	for _, collName := range requiredCollections {
		if !existingCollections[collName] {
			log.Printf("Creating collection: %s", collName)
			err := db.CreateCollection(ctx, collName)
			if err != nil {
				return fmt.Errorf("failed to create collection %s: %w", collName, err)
			}
			log.Printf("Successfully created collection: %s", collName)
		} else {
			log.Printf("Collection already exists: %s", collName)
		}
	}

	return nil
}

// GetAllPosts retrieves all posts from the posts collection
func GetAllPosts(ctx context.Context, db *mongo.Database) ([]models.Post, error) {
	collection := db.Collection("posts")

	// Find all posts
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find posts: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode all posts
	var posts []models.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, fmt.Errorf("failed to decode posts: %w", err)
	}

	return posts, nil
}
