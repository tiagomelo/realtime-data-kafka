package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"github.com/tiagomelo/realtime-data-kafka/config"
	"github.com/tiagomelo/realtime-data-kafka/screen"
	"github.com/tiagomelo/realtime-data-kafka/stats"
)

const bootstrapServersKey = "bootstrap.servers"

func stringPrt(s string) *string {
	return &s
}

func run(log *log.Logger, cfg *config.Config, transactionsFile string) error {
	log.Println("main: Initializing Kafka producer")
	defer log.Println("main: Completed")
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		bootstrapServersKey: cfg.KafkaBrokerHost,
	})
	if err != nil {
		return errors.Wrap(err, "creating producer")
	}
	defer producer.Close()
	file, err := os.Open(transactionsFile)
	if err != nil {
		return errors.Wrapf(err, "opening file %s", transactionsFile)
	}
	defer file.Close()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	stats := &stats.KafkaProducerStats{}
	screen, err := screen.NewKafkaProducerScreen(stats)
	if err != nil {
		return errors.New("starting screen")
	}

	start := time.Now()

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(1))
			stats.UpdateElapsedTime(time.Since(start))
			screen.UpdateContent(false)
		}
	}()

	deliveryChan := make(chan kafka.Event)
	scanner := bufio.NewScanner(file)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			if err := producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: stringPrt(cfg.KafkaTopic), Partition: kafka.PartitionAny},
				Value:          []byte(line),
			}, deliveryChan); err != nil {
				log.Printf("%v when publishing to kafka topic %s", err, cfg.KafkaTopic)
			}
			stats.IncrTotalPublishedMessages()
			delivery := <-deliveryChan
			m := delivery.(*kafka.Message)
			if m.TopicPartition.Error != nil {
				stats.IncrTotalFailedMessageDeliveries()
			}
		}
		if err := scanner.Err(); err != nil {
			errors.Wrapf(err, "reading file %s", transactionsFile)
		}
	}()

	// Wait for any error or interrupt signal.
	select {
	case err := <-serverErrors:
		return err
	case sig := <-shutdown:
		screen.UpdateContent(true)
		log.Printf("run: %v: Start shutdown", sig)
		return nil
	}
}

var opts struct {
	File string `short:"f" long:"file" description:"input file" required:"true"`
}

func main() {
	const (
		envFile     = ".env"
		logFileName = "logs/producer.txt"
	)
	flags.ParseArgs(&opts, os.Args)
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf(`opening log file "%s": %v`, logFileName, err)
	}
	log := log.New(logFile, "KAFKA PRODUCER : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	cfg, err := config.Read(envFile)
	if err != nil {
		log.Println(errors.Wrap(err, "reading config"))
		fmt.Println(errors.Wrap(err, "reading config"))
		os.Exit(1)
	}
	if err := run(log, cfg, opts.File); err != nil {
		log.Println(err)
		fmt.Println(err)
		os.Exit(1)
	}
}
