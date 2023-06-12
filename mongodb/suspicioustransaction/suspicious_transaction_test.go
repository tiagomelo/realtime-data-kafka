// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package suspicioustransaction

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tiagomelo/realtime-data-kafka/mongodb"
	"github.com/tiagomelo/realtime-data-kafka/mongodb/suspicioustransaction/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestInsert(t *testing.T) {
	testCases := []struct {
		name                     string
		mockCollection           func(mongoClient *mongo.Client, databaseName, collectionName string) *mongo.Collection
		mockInsertIntoCollection func(ctx context.Context, collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error)
		expectedError            error
	}{
		{
			name: "happy path",
			mockCollection: func(mongoClient *mongo.Client, databaseName, collectionName string) *mongo.Collection {
				return new(mongo.Collection)
			},
			mockInsertIntoCollection: func(ctx context.Context, collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
				return new(mongo.InsertOneResult), nil
			},
		},
		{
			name: "error",
			mockCollection: func(mongoClient *mongo.Client, databaseName, collectionName string) *mongo.Collection {
				return new(mongo.Collection)
			},
			mockInsertIntoCollection: func(ctx context.Context, collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
				return nil, errors.New("random error")
			},
			expectedError: errors.New("inserting suspicious transaction: random error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			collection = tc.mockCollection
			insertIntoCollection = tc.mockInsertIntoCollection
			err := Insert(context.TODO(), new(mongodb.MongoDb), new(models.SuspiciousTransaction))
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.expectedError)
				}
			}
		})
	}
}
