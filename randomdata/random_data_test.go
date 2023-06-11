package randomdata

import (
	"testing"
	"time"
)

func TestTransactionID(t *testing.T) {
	id := TransactionID()
	if id < 1111111111 || id > 9999999999 {
		t.Errorf("Transaction ID out of range: %d", id)
	}
}

func TestAccountNumber(t *testing.T) {
	accountNumber := AccountNumber()
	if accountNumber < 111111111 || accountNumber > 999999999 {
		t.Errorf("Account number out of range: %d", accountNumber)
	}
}

func TestTransactionAmount(t *testing.T) {
	minAmount := 100
	maxAmount := 1000
	amount := TransactionAmount(float32(minAmount), float32(maxAmount))
	if amount < float32(minAmount) || amount > float32(maxAmount) {
		t.Errorf("Transaction amount out of range: %f", amount)
	}
}

func TestTransactionTime(t *testing.T) {
	result := TransactionTime()
	currentTime := time.Now()
	maxTime := currentTime.Add(-24 * time.Hour)
	if result.Before(maxTime) || result.After(currentTime) {
		t.Errorf("Returned time is not within the acceptable range: %v", result)
	}
}

func TestLocation(t *testing.T) {
	location := Location()
	found := false
	for _, loc := range locations {
		if location == loc {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Invalid random location: %s", location)
	}
}
