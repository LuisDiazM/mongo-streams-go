package main

import (
	"github.com/LuisDiazM/mongo-streams-go/app"
	"github.com/LuisDiazM/mongo-streams-go/infraestructure/database"
)

func main() {

	//manual dependency injection
	settings := app.GetAppSettings()
	mongoImp := database.NewMongoImplmentation(settings.MongoDb)
	application := app.NewApplication(mongoImp, settings)

	application.Start()
}
