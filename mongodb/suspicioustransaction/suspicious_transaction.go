package suspicioustransaction

import (
	"context"

	"github.com/pkg/errors"
	"github.com/tiagomelo/realtime-data-kafka/mongodb"
	"github.com/tiagomelo/realtime-data-kafka/mongodb/suspicioustransaction/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// For ease of unit testing.
var (
	collection = func(mongoClient *mongo.Client, databaseName, collectionName string) *mongo.Collection {
		return mongoClient.Database(databaseName).Collection(collectionName)
	}
	insertIntoCollection = func(ctx context.Context, collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
		return collection.InsertOne(ctx, document)
	}
)

// Insert inserts a new suspicious transaction into the MongoDB collection.
func Insert(ctx context.Context, db *mongodb.MongoDb, newSuspiciousTran *models.SuspiciousTransaction) error {
	const collectionName = "suspicious_transactions"
	coll := collection(db.Client, db.DatabaseName, collectionName)
	_, err := insertIntoCollection(ctx, coll, newSuspiciousTran)
	if err != nil {
		return errors.Wrap(err, "inserting suspicious transaction")
	}
	return nil
}
