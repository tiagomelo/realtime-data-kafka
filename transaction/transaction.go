// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package transaction

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

// Transaction represents a transaction message.
type Transaction struct {
	TransactionID     int       `json:"transaction_id"`
	AccountNumber     int       `json:"account_number"`
	TransactionType   string    `json:"transaction_type"`
	TransactionAmount float32   `json:"transaction_amount"`
	TransactionTime   time.Time `json:"transaction_time"`
	Location          string    `json:"location"`
}

// New creates a new Transaction from the raw JSON transaction data.
func New(rawTransaction string) (*Transaction, error) {
	t := new(Transaction)
	if err := json.Unmarshal([]byte(rawTransaction), &t); err != nil {
		return nil, errors.Wrap(err, "unmarshalling transaction")
	}
	return t, nil
}

// IsSuspicious checks if the transaction amount is suspicious.
func (t *Transaction) IsSuspicious() bool {
	const suspiciousAmount = float32(10_000)
	return t.TransactionAmount > suspiciousAmount
}
