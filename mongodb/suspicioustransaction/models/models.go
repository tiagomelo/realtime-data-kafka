package models

import "time"

// SuspiciousTransaction represents a suspicious transaction.
type SuspiciousTransaction struct {
	TransactionId     int       `bson:"transaction_id"`
	AccountNumber     int       `bson:"account_number"`
	TransactionType   string    `bson:"transaction_type"`
	TransactionAmount float32   `bson:"transaction_amount"`
	TransactionTime   time.Time `bson:"transaction_time"`
	Location          string    `bson:"location"`
}
