package config

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type AppConfig struct {
	Version                        string
	Port                           uint16
	MongoDbURI                     string
	IcsFileScrapeIntervalInSeconds int64
}

var Config *AppConfig
var MongoDB *mongo.Database
