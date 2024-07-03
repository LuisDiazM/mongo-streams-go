package app

import (
	"log"

	"github.com/LuisDiazM/mongo-streams-go/infraestructure/database"
	"github.com/kelseyhightower/envconfig"
)

type AppSettings struct {
	MongoDb *database.MongoSettings
}

func GetAppSettings() *AppSettings {
	var s AppSettings
	err := envconfig.Process("backups", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &s
}

func GetMongoSettings() *database.MongoSettings {
	return GetAppSettings().MongoDb
}
