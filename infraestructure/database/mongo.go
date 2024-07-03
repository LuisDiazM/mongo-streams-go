package database

import (
	"context"
	"log"
	"time"

	"github.com/LuisDiazM/mongo-streams-go/infraestructure/database/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	watchInsert = "insert"
	watchUpdate = "update"
)

type MongoImp struct {
	Client   *mongo.Client
	settings *MongoSettings
	ctx      context.Context
}

func NewMongoImplmentation(settings *MongoSettings) *MongoImp {
	ctx := context.TODO()
	clientOptions := options.Client().ApplyURI(settings.Url).SetMaxPoolSize(settings.MaxPoolSize)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return &MongoImp{settings: settings, ctx: ctx, Client: client}
}

func (db *MongoImp) Ping() {
	err := db.Client.Ping(db.ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!!")
}

func (db *MongoImp) WatchCollection() {
	database := db.Client.Database(db.settings.Database)
	collection := database.Collection(db.settings.BaseCollection)
	newCollection := database.Collection(db.settings.NewCollection)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "operationType", Value: bson.D{
				{Key: "$in", Value: bson.A{watchInsert, watchUpdate}},
			}},
		}}},
	}
	opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	changeStream, err := collection.Watch(db.ctx, pipeline, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer changeStream.Close(db.ctx)

	for changeStream.Next(db.ctx) {
		var event map[string]interface{}
		if err := changeStream.Decode(&event); err != nil {
			log.Fatal(err)
		}
		db.handleEvent(newCollection, event, db.ctx)
	}
	if err := changeStream.Err(); err != nil {
		log.Fatal(err)
	}
}

func (db *MongoImp) handleEvent(newCollection *mongo.Collection, event map[string]interface{}, ctx context.Context) {
	t1 := time.Now()
	operationType, operationTypeOK := event["operationType"].(string)
	if !operationTypeOK {
		return
	}
	defer utils.ProcessingTime(t1, operationType)
	fullDocument, fullDocumentOK := event["fullDocument"].(map[string]interface{})
	switch operationType {
	case "insert":
		if fullDocumentOK {
			_, err := newCollection.InsertOne(ctx, fullDocument)
			if err != nil {
				log.Printf("Error occurred: %v", err)
				return
			}
		}
	case "update":
		documentKey, documentKeyOK := event["documentKey"].(map[string]interface{})
		id, idOK := documentKey["_id"].(primitive.ObjectID)
		if !documentKeyOK || !idOK {
			return
		}
		updateDescription, updateDescriptionOK := event["updateDescription"].(map[string]interface{})
		removeFields, removeFieldsOK := updateDescription["removedFields"].(primitive.A)

		if !updateDescriptionOK || !removeFieldsOK {
			return
		}

		if len(removeFields) > 0 {
			filedsToRemove := utils.MapFields(removeFields)
			_, err := newCollection.UpdateOne(ctx, bson.D{{Key: "_id", Value: id}}, filedsToRemove)
			if err != nil {
				log.Printf("Error occurred: %v", err)
				return
			}
			return
		} else {
			if fullDocumentOK {
				_, err := newCollection.UpdateOne(ctx, bson.D{{Key: "_id", Value: id}}, bson.D{{Key: "$set", Value: fullDocument}})
				if err != nil {
					log.Printf("Error occurred: %v", err)
					return
				}
			}
		}
	}
}
