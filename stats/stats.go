// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package stats

import (
	"sync/atomic"
	"time"
)

// KafkaConsumerStats represents the statistics for Kafka consumer operations.
type KafkaConsumerStats struct {
	totalTransactions                      int64
	totalSuspiciousTransactions            int64
	totalUnmarshallingMsgErrors            int64
	totalInsertSuspiciousTransactionErrors int64
	elapsedTime                            time.Duration
}

// KafkaProducerStats represents the statistics for Kafka producer operations.
type KafkaProducerStats struct {
	totalPublishedMessages       int64
	totalFailedMessageDeliveries int64
	elapsedTime                  time.Duration
}

// IncrTotalTransactions increments the total number of transactions.
func (stats *KafkaConsumerStats) IncrTotalTransactions() {
	atomic.AddInt64(&stats.totalTransactions, 1)
}

// TotalTransactions returns the total number of transactions.
func (stats *KafkaConsumerStats) TotalTransactions() int64 {
	return stats.totalTransactions
}

// IncrTotalSuspiciousTransactions increments the total number of suspicious transactions.
func (stats *KafkaConsumerStats) IncrTotalSuspiciousTransactions() {
	atomic.AddInt64(&stats.totalSuspiciousTransactions, 1)
}

// TotalSuspiciousTransactions returns the total number of suspicious transactions.
func (stats *KafkaConsumerStats) TotalSuspiciousTransactions() int64 {
	return stats.totalSuspiciousTransactions
}

// IncrTotalUnmarshallingMsgErrors increments the total number of unmarshalling message errors.
func (stats *KafkaConsumerStats) IncrTotalUnmarshallingMsgErrors() {
	atomic.AddInt64(&stats.totalUnmarshallingMsgErrors, 1)
}

// TotalUnmarshallingMsgErrors returns the total number of unmarshalling message errors.
func (stats *KafkaConsumerStats) TotalUnmarshallingMsgErrors() int64 {
	return stats.totalUnmarshallingMsgErrors
}

// IncrTotalInsertSuspiciousTransactionErrors increments the total number of insert suspicious transaction errors.
func (stats *KafkaConsumerStats) IncrTotalInsertSuspiciousTransactionErrors() {
	atomic.AddInt64(&stats.totalInsertSuspiciousTransactionErrors, 1)
}

// TotalInsertSuspiciousTransactionErrors returns the total number of insert suspicious transaction errors.
func (stats *KafkaConsumerStats) TotalInsertSuspiciousTransactionErrors() int64 {
	return stats.totalInsertSuspiciousTransactionErrors
}

// UpdateElapsedTime updates the elapsed time for Kafka consumer operations.
func (stats *KafkaConsumerStats) UpdateElapsedTime(elapsedTime time.Duration) {
	stats.elapsedTime = elapsedTime
}

// ElapsedTime returns the elapsed time for Kafka consumer operations.
func (stats *KafkaConsumerStats) ElapsedTime() time.Duration {
	return stats.elapsedTime
}

// IncrTotalPublishedMessages increments the total number of published messages.
func (stats *KafkaProducerStats) IncrTotalPublishedMessages() {
	atomic.AddInt64(&stats.totalPublishedMessages, 1)
}

// TotalPublishedMessages returns the total number of published messages.
func (stats *KafkaProducerStats) TotalPublishedMessages() int64 {
	return stats.totalPublishedMessages
}

// IncrTotalFailedMessageDeliveries increments the total number of failed message deliveries.
func (stats *KafkaProducerStats) IncrTotalFailedMessageDeliveries() {
	atomic.AddInt64(&stats.totalFailedMessageDeliveries, 1)
}

// TotalFailedMessageDeliveries returns the total number of failed message deliveries.
func (stats *KafkaProducerStats) TotalFailedMessageDeliveries() int64 {
	return stats.totalFailedMessageDeliveries
}

// UpdateElapsedTime updates the elapsed time for Kafka producer operations.
func (stats *KafkaProducerStats) UpdateElapsedTime(elapsedTime time.Duration) {
	stats.elapsedTime = elapsedTime
}

// ElapsedTime returns the elapsed time for Kafka producer operations.
func (stats *KafkaProducerStats) ElapsedTime() time.Duration {
	return stats.elapsedTime
}
