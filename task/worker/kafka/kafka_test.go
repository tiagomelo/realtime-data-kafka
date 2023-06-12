// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package kafka

import (
	"context"
	"errors"
	"log"
	"strings"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/require"
	"github.com/tiagomelo/realtime-data-kafka/mongodb"
	"github.com/tiagomelo/realtime-data-kafka/mongodb/suspicioustransaction/models"
	"github.com/tiagomelo/realtime-data-kafka/stats"
	"github.com/tiagomelo/realtime-data-kafka/stringify"
)

func TestWork(t *testing.T) {
	testCases := []struct {
		name                                           string
		msg                                            string
		mockPrintToLog                                 func(log *log.Logger, v ...any)
		mockStInsert                                   func(ctx context.Context, db *mongodb.MongoDb, sp *models.SuspiciousTransaction) error
		expectedTotalTransactions                      int64
		expectedTotalUnmarshallingMsgErrors            int64
		expectedTotalSuspiciousTransactions            int64
		expectedTotalInsertSuspiciousTransactionErrors int64
	}{
		{
			name:                      "no suspicious transactions",
			msg:                       `{"transaction_id":5699757367,"account_number":215489034,"transaction_type":"withdrawal","transaction_amount":1308.58,"transaction_time":"2023-06-05T03:05:12.495058-03:00","location":"Fort Worth, TX"}`,
			mockPrintToLog:            func(log *log.Logger, v ...any) {},
			expectedTotalTransactions: int64(1),
		},
		{
			name: "with suspicious transaction",
			msg:  `{"transaction_id":5699757367,"account_number":215489034,"transaction_type":"withdrawal","transaction_amount":11308.58,"transaction_time":"2023-06-05T03:05:12.495058-03:00","location":"Fort Worth, TX"}`,
			mockPrintToLog: func(log *log.Logger, v ...any) {
				m := stringify.VariadicToStringArray(v)
				var contains bool
				for _, v := range m {
					if strings.Contains(v, "suspicious transaction:") {
						contains = true
					}
				}
				require.True(t, contains)
			},
			mockStInsert: func(ctx context.Context, db *mongodb.MongoDb, sp *models.SuspiciousTransaction) error {
				return nil
			},
			expectedTotalTransactions:           int64(1),
			expectedTotalSuspiciousTransactions: int64(1),
		},
		{
			name: "invalid message",
			msg:  "blabla",
			mockPrintToLog: func(log *log.Logger, v ...any) {
				expectedMsg := []string{"[checking if transaction is suspicious: unmarshalling transaction: invalid character 'b' looking for beginning of value]"}
				c := stringify.VariadicToStringArray(v)
				require.Equal(t, expectedMsg, c)
			},
			expectedTotalTransactions:           int64(1),
			expectedTotalUnmarshallingMsgErrors: int64(1),
		},
		{
			name: "error when saving suspicious transaction to db",
			msg:  `{"transaction_id":5699757367,"account_number":215489034,"transaction_type":"withdrawal","transaction_amount":11308.58,"transaction_time":"2023-06-05T03:05:12.495058-03:00","location":"Fort Worth, TX"}`,
			mockPrintToLog: func(log *log.Logger, v ...any) {
				var contains bool
				m := stringify.VariadicToStringArray(v)
				for _, v := range m {
					if !strings.Contains(v, "suspicious transaction:") {
						if strings.Contains(v, "error when inserting") {
							contains = true
						}
						require.True(t, contains)
					}
				}
			},
			mockStInsert: func(ctx context.Context, db *mongodb.MongoDb, sp *models.SuspiciousTransaction) error {
				return errors.New("random error")
			},
			expectedTotalTransactions:                      int64(1),
			expectedTotalSuspiciousTransactions:            int64(1),
			expectedTotalInsertSuspiciousTransactionErrors: int64(1),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stats := new(stats.KafkaConsumerStats)
			printToLog = tc.mockPrintToLog
			stInsert = tc.mockStInsert
			worker := &Worker{
				Stats: stats,
				Msg: &kafka.Message{
					Value: []byte(tc.msg),
				},
			}
			worker.Work(context.TODO())
			require.Equal(t, tc.expectedTotalTransactions, stats.TotalTransactions())
			require.Equal(t, tc.expectedTotalSuspiciousTransactions, stats.TotalSuspiciousTransactions())
			require.Equal(t, tc.expectedTotalInsertSuspiciousTransactionErrors, stats.TotalInsertSuspiciousTransactionErrors())
		})
	}
}
