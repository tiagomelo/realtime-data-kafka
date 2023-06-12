// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package mongodb

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestConnect(t *testing.T) {
	testCases := []struct {
		name          string
		mockNewClient func(opts ...*options.ClientOptions) (*mongo.Client, error)
		mockConnect   func(ctx context.Context, client *mongo.Client) error
		mockPing      func(ctx context.Context, client *mongo.Client) error
		expectedError error
	}{
		{
			name: "happy path",
			mockNewClient: func(opts ...*options.ClientOptions) (*mongo.Client, error) {
				return new(mongo.Client), nil
			},
			mockConnect: func(ctx context.Context, client *mongo.Client) error {
				return nil
			},
			mockPing: func(ctx context.Context, client *mongo.Client) error {
				return nil
			},
		},
		{
			name: "error when creating client",
			mockNewClient: func(opts ...*options.ClientOptions) (*mongo.Client, error) {
				return nil, errors.New("random error")
			},
			expectedError: errors.New("failed to create MongoDB client: random error"),
		},
		{
			name: "error when connecting",
			mockNewClient: func(opts ...*options.ClientOptions) (*mongo.Client, error) {
				return new(mongo.Client), nil
			},
			mockConnect: func(ctx context.Context, client *mongo.Client) error {
				return errors.New("random error")
			},
			expectedError: errors.New("failed to connect to MongoDB server: random error"),
		},
		{
			name: "error when doing ping",
			mockNewClient: func(opts ...*options.ClientOptions) (*mongo.Client, error) {
				return new(mongo.Client), nil
			},
			mockConnect: func(ctx context.Context, client *mongo.Client) error {
				return nil
			},
			mockPing: func(ctx context.Context, client *mongo.Client) error {
				return errors.New("random error")
			},
			expectedError: errors.New("failed to ping MongoDB server: random error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newClient = tc.mockNewClient
			connect = tc.mockConnect
			ping = tc.mockPing
			conn, err := Connect(context.TODO(), "host", "database", 111)
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.expectedError)
				}
				require.NotNil(t, conn)
			}
		})
	}
}
