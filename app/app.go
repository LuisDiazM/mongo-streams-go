package app

import (
	"github.com/LuisDiazM/mongo-streams-go/infraestructure/database"
)

type Application struct {
	MongoImp *database.MongoImp
	Settings *AppSettings
}

func NewApplication(mongo *database.MongoImp, settings *AppSettings) *Application {
	return &Application{MongoImp: mongo, Settings: settings}
}

func (appl *Application) Start() {
	appl.MongoImp.Ping()
	appl.MongoImp.WatchCollection()
}
