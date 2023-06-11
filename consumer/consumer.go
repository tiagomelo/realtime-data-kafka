package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
	"github.com/tiagomelo/realtime-data-kafka/config"
	"github.com/tiagomelo/realtime-data-kafka/mongodb"
	"github.com/tiagomelo/realtime-data-kafka/screen"
	"github.com/tiagomelo/realtime-data-kafka/stats"
	"github.com/tiagomelo/realtime-data-kafka/task"
	kafkaWorker "github.com/tiagomelo/realtime-data-kafka/task/worker/kafka"
)

// Useful constants.
const (
	bootstrapServersKey   = "bootstrap.servers"
	groupIdKey            = "group.id"
	autoOffsetResetKey    = "auto.offset.reset"
	autoOffsetReset       = "earliest"
	enablePartitionEofKey = "enable.partition.eof"
)

func run(log *log.Logger) error {
	const envFile = ".env"
	log.Println("main: Initializing Kafka consumer")
	defer log.Println("main: Completed")
	ctx := context.Background()

	cfg, err := config.Read(envFile)
	if err != nil {
		return errors.Wrap(err, "reading config")
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		bootstrapServersKey:   cfg.KafkaBrokerHost,
		groupIdKey:            cfg.KafkaGroupId,
		autoOffsetResetKey:    autoOffsetReset,
		enablePartitionEofKey: false,
	})
	if err != nil {
		return errors.Wrapf(err, "connecting to broker %s", cfg.KafkaBrokerHost)
	}

	if err := consumer.SubscribeTopics([]string{cfg.KafkaTopic}, nil); err != nil {
		return errors.Wrapf(err, "subscribing to topic %s", cfg.KafkaTopic)
	}

	db, err := mongodb.Connect(ctx, cfg.MongodbHostName, cfg.MongodbDatabase, cfg.MongodbPort)
	if err != nil {
		return errors.Wrapf(err, "connecting to mongodb")
	}

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	maxGoRoutines := runtime.GOMAXPROCS(0)
	pool := task.New(ctx, maxGoRoutines)

	stats := &stats.KafkaConsumerStats{}
	screen, err := screen.NewKafkaConsumerScreen(stats)
	if err != nil {
		return errors.New("starting screen")
	}

	start := time.Now()

	go func() {
		defer close(shutdown)
		defer close(serverErrors)
		for {
			select {
			case <-shutdown:
				log.Printf("run: Start shutdown")
				if err := consumer.Close(); err != nil {
					serverErrors <- errors.Wrap(err, "closing Kafka consumer")
				}
				return
			default:
				msg, err := consumer.ReadMessage(-1)
				if err != nil {
					serverErrors <- err
				} else {
					kw := &kafkaWorker.Worker{Msg: msg, Stats: stats, Db: db, Log: log}
					pool.Do(kw)
				}
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(1))
			stats.UpdateElapsedTime(time.Since(start))
			screen.UpdateContent(false)
		}
	}()

	// Wait for any error or interrupt signal.
	select {
	case err := <-serverErrors:
		return err
	case sig := <-shutdown:
		screen.UpdateContent(true)
		log.Printf("run: %v: Start shutdown", sig)
		// Asking listener to shutdown and shed load.
		if err := consumer.Close(); err != nil {
			return errors.Wrap(err, "closing Kafka consumer")
		}
		return nil
	}
}

func main() {
	const logFileName = "logs/consumer.txt"
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf(`opening log file "%s": %v`, logFileName, err)
	}
	log := log.New(logFile, "KAFKA CONSUMER : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	if err := run(log); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
