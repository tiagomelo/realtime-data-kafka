package transaction

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	timestamp := "2023-06-05T03:05:12.495058-03:00"
	parsedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name           string
		input          string
		expectedOutput *Transaction
		expectedError  error
	}{
		{
			name:  "happy path",
			input: `{"transaction_id":5699757367,"account_number":215489034,"transaction_type":"withdrawal","transaction_amount":11308.58,"transaction_time":"2023-06-05T03:05:12.495058-03:00","location":"Fort Worth, TX"}`,
			expectedOutput: &Transaction{
				TransactionID:     5699757367,
				AccountNumber:     215489034,
				TransactionType:   "withdrawal",
				TransactionAmount: 11308.58,
				TransactionTime:   parsedTime,
				Location:          "Fort Worth, TX",
			},
		},
		{
			name:          "error",
			input:         `invalid input`,
			expectedError: errors.New("unmarshalling transaction: invalid character 'i' looking for beginning of value"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := New(tc.input)
			if err != nil {
				if tc.expectedError == nil {
					t.Fatalf(`expected no error, got "%v"`, err)
				}
				require.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				if tc.expectedError != nil {
					t.Fatalf(`expected error "%v", got nil`, tc.expectedError)
				}
				require.Equal(t, tc.expectedOutput, output)
			}
		})
	}
}

func TestIsSuspicious(t *testing.T) {
	testCases := []struct {
		name           string
		input          *Transaction
		expectedOutput bool
	}{
		{
			name: "suspicious",
			input: &Transaction{
				TransactionAmount: 11308.58,
			},
			expectedOutput: true,
		},
		{
			name: "not suspicious",
			input: &Transaction{
				TransactionAmount: 1308.58,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			suspicious := tc.input.IsSuspicious()
			require.Equal(t, tc.expectedOutput, suspicious)
		})
	}
}
