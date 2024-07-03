package database

type MongoSettings struct {
	Url            string `envconfig:"MONGO_URL" required:"true"`
	MaxPoolSize    uint64 `envconfig:"MONGO_MAXPOOLSIZE" default:"20"`
	Database       string `envconfig:"MONGO_DATABASE" required:"true"`
	BaseCollection string `envconfig:"MONGO_BASE_COLLECTION" required:"true"`
	NewCollection  string `envconfig:"MONGO_NEW_COLLECTION" required:"true"`
}
