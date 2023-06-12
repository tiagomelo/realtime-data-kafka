// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/tiagomelo/realtime-data-kafka/mongodb"
	"github.com/tiagomelo/realtime-data-kafka/mongodb/suspicioustransaction"
	"github.com/tiagomelo/realtime-data-kafka/mongodb/suspicioustransaction/models"
	"github.com/tiagomelo/realtime-data-kafka/stats"
	"github.com/tiagomelo/realtime-data-kafka/transaction"
)

// For ease of unit testing.
var (
	printToLog = func(log *log.Logger, v ...any) {
		log.Println(v...)
	}
	stInsert = func(ctx context.Context, db *mongodb.MongoDb, sp *models.SuspiciousTransaction) error {
		return suspicioustransaction.Insert(ctx, db, sp)
	}
)

// Worker represents a Kafka consumer worker.
type Worker struct {
	Msg   *kafka.Message
	Stats *stats.KafkaConsumerStats
	Db    *mongodb.MongoDb
	Log   *log.Logger
}

// insertSuspiciousTransaction inserts a suspicious transaction into MongoDB.
func (c *Worker) insertSuspiciousTransaction(ctx context.Context, sp *transaction.Transaction) error {
	spDb := &models.SuspiciousTransaction{
		TransactionId:     sp.TransactionID,
		AccountNumber:     sp.AccountNumber,
		TransactionType:   sp.TransactionType,
		TransactionAmount: sp.TransactionAmount,
		TransactionTime:   sp.TransactionTime,
		Location:          sp.Location,
	}
	return stInsert(ctx, c.Db, spDb)
}

// Work processes the Kafka message and performs the necessary operations.
func (c *Worker) Work(ctx context.Context) {
	c.Stats.IncrTotalTransactions()
	transaction, err := transaction.New(string(c.Msg.Value))
	if err != nil {
		c.Stats.IncrTotalUnmarshallingMsgErrors()
		printToLog(c.Log, fmt.Errorf("checking if transaction is suspicious: %v", err))
		return
	}
	if transaction.IsSuspicious() {
		c.Stats.IncrTotalSuspiciousTransactions()
		printToLog(c.Log, "suspicious transaction: %+v\n", transaction)
		if err := c.insertSuspiciousTransaction(ctx, transaction); err != nil {
			c.Stats.IncrTotalInsertSuspiciousTransactionErrors()
			printToLog(c.Log, "error when inserting suspicious transaction in mongodb %+v: %v\n", transaction, err)
		}
	}
}
