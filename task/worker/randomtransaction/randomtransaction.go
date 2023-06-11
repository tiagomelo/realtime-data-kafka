package randomtransaction

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/tiagomelo/realtime-data-kafka/randomdata"
	"github.com/tiagomelo/realtime-data-kafka/transaction"
)

// For ease of unit testing.
var (
	openFile        = os.OpenFile
	jsonMarshal     = json.Marshal
	fileWriteString = func(file *os.File, s string) (n int, err error) {
		return file.WriteString(s)
	}
	printToLog = func(log *log.Logger, v ...any) {
		log.Println(v...)
	}
)

// Worker generates random transaction data.
type Worker struct {
	FilePath  string
	MinAmount float32
	MaxAmount float32
	Log       *log.Logger
}

// Work generates a random transaction and writes it to a file.
func (w *Worker) Work(ctx context.Context) {
	t := generateRandomTransaction(w.MinAmount, w.MaxAmount)
	file, err := openFile(w.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		printToLog(w.Log, "error opening file:", err)
		return
	}
	defer file.Close()
	jsonData, err := jsonMarshal(t)
	if err != nil {
		printToLog(w.Log, "error marshalling json:", err)
		return
	}
	_, err = fileWriteString(file, string(jsonData)+"\n")
	if err != nil {
		printToLog(w.Log, "error writing to file:", err)
	}
}

// generateRandomTransaction generates a random transaction with the given minimum and maximum amounts.
func generateRandomTransaction(minAmount, maxAmount float32) *transaction.Transaction {
	const withdrawal = "withdrawal"
	t := &transaction.Transaction{
		TransactionID:     randomdata.TransactionID(),
		AccountNumber:     randomdata.AccountNumber(),
		TransactionType:   withdrawal,
		TransactionAmount: randomdata.TransactionAmount(minAmount, maxAmount),
		TransactionTime:   randomdata.TransactionTime(),
		Location:          randomdata.Location(),
	}
	return t
}
