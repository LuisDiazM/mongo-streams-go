package utils

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapFields(removeFields primitive.A) bson.D {
	update := bson.D{{Key: "$unset", Value: bson.D{}}}
	for _, field := range removeFields {
		update[0].Value = append(update[0].Value.(bson.D), bson.E{Key: field.(string), Value: ""})
	}
	return update
}
func ProcessingTime(t1 time.Time, operation string) {
	t2 := time.Since(t1).Milliseconds()
	log.Printf(`%s took %d ms \n`, operation, t2)
}
