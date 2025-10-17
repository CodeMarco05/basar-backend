package config

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type AppConfig struct {
	Version string
	Port    uint16
}

var Config *AppConfig
var MongoDB *mongo.Database
