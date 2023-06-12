// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package randomdata

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// locations is a slice of pre-defined locations for generating random transaction locations.
var locations = []string{
	"New York, NY",
	"Los Angeles, CA",
	"Chicago, IL",
	"Houston, TX",
	"Phoenix, AZ",
	"Philadelphia, PA",
	"San Antonio, TX",
	"San Diego, CA",
	"Dallas, TX",
	"San Jose, CA",
	"Austin, TX",
	"Jacksonville, FL",
	"Fort Worth, TX",
	"Columbus, OH",
	"Charlotte, NC",
	"San Francisco, CA",
	"Indianapolis, IN",
	"Seattle, WA",
	"Denver, CO",
	"Washington, DC",
}

// TransactionID generates a random transaction ID.
func TransactionID() int {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	r.Seed(time.Now().UnixNano())
	return r.Intn(9999999999-1111111111+1) + 1111111111
}

// AccountNumber generates a random account number.
func AccountNumber() int {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	r.Seed(time.Now().UnixNano())
	return r.Intn(999999999-111111111+1) + 111111111
}

// TransactionAmount generates a random transaction amount between the specified minimum and maximum amounts.
func TransactionAmount(minAmount, maxAmount float32) float32 {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	randomAmount := r.Float32()*(maxAmount-minAmount) + minAmount
	formattedAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", randomAmount), 32)
	return float32(formattedAmount)
}

// TransactionTime generates a random transaction time within the last 24 hours.
func TransactionTime() time.Time {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	randomDuration := time.Duration(r.Intn(86400)) * time.Second
	randomTime := time.Now().Add(-randomDuration)
	return randomTime
}

// Location generates a random transaction location from the pre-defined locations.
func Location() string {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	return locations[r.Intn(len(locations))]
}
