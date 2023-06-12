// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
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
