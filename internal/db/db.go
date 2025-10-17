package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBClient defines the interface for MongoDB operations
type MongoDBClient interface {
	Init(uri string, database string) error
	HealthCheck(ctx context.Context) error
	GetClient() *mongo.Client
	GetDatabase() *mongo.Database
	Close(ctx context.Context) error
}

// mongoClient implements the MongoDBClient interface
type mongoClient struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewMongoDBClient creates a new MongoDB client instance
func NewMongoDBClient() MongoDBClient {
	return &mongoClient{}
}

// Init initializes the MongoDB client with connection pooling and retry logic
func (m *mongoClient) Init(uri string, databaseName string) error {
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

	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("Attempting to connect to MongoDB (attempt %d/%d)...", attempt, maxRetries)

		m.client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Printf("Failed to connect to MongoDB (attempt %d): %v", attempt, err)
			if attempt < maxRetries {
				time.Sleep(retryDelay)
				continue
			}
			return fmt.Errorf("failed to connect to MongoDB after %d attempts: %w", maxRetries, err)
		}

		// Check the connection
		err = m.client.Ping(context.TODO(), nil)
		if err != nil {
			log.Printf("Failed to ping MongoDB (attempt %d): %v", attempt, err)
			if attempt < maxRetries {
				time.Sleep(retryDelay)
				continue
			}
			return fmt.Errorf("failed to ping MongoDB after %d attempts: %w", maxRetries, err)
		}

		log.Println("Successfully connected to MongoDB!")
		break
	}

	// Get a handle for the database
	m.database = m.client.Database(databaseName)

	// Create collections if they don't exist
	if err := m.ensureCollectionsExist(context.TODO()); err != nil {
		return fmt.Errorf("failed to ensure collections exist: %w", err)
	}

	return nil
}

// HealthCheck performs a health check on the MongoDB connection
func (m *mongoClient) HealthCheck(ctx context.Context) error {
	if m.client == nil {
		return fmt.Errorf("MongoDB client is not initialized")
	}

	// Ping the MongoDB server
	if err := m.client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("MongoDB health check failed: %w", err)
	}

	return nil
}

// GetClient returns the underlying MongoDB client
func (m *mongoClient) GetClient() *mongo.Client {
	return m.client
}

// GetDatabase returns the database instance
func (m *mongoClient) GetDatabase() *mongo.Database {
	return m.database
}

// Close closes the MongoDB connection
func (m *mongoClient) Close(ctx context.Context) error {
	if m.client == nil {
		return nil
	}

	if err := m.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	log.Println("MongoDB connection closed")
	return nil
}

// ensureCollectionsExist creates required collections if they don't already exist
func (m *mongoClient) ensureCollectionsExist(ctx context.Context) error {
	// List existing collections
	collections, err := m.database.ListCollectionNames(ctx, bson.M{})
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
		"listings",
		// Add more collection names as needed
		// "users",
		// etc.
	}

	// Create missing collections
	for _, collName := range requiredCollections {
		if !existingCollections[collName] {
			log.Printf("Creating collection: %s", collName)
			err := m.database.CreateCollection(ctx, collName)
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
