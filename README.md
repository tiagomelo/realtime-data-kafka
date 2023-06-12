# realtime-data-kafka

This is a tutorial showing real time processing using [Kafka](https://kafka.apache.org/) and [MongoDB](https://www.mongodb.com/).

Article: [Real time data processing: easily processing 10 million messages with Golang, Kafka and MongoDB](https://www.linkedin.com/pulse/real-time-data-processing-easily-10-million-messages-golang-melo)

## requirements
- [Kafka](https://kafka.apache.org/)
- [MongoDB](https://www.mongodb.com/)

## before running both producer and consumer

```
make zookeeper
make kafka
```

## producer

It accepts a file with fictitous transactions in JSON format and publish every line into a Kafka topic.

JSON example:

```
{
  "transaction_id": 4508561159,
  "account_number": 395402066,
  "transaction_type": "withdrawal",
  "transaction_amount": 2718.79,
  "transaction_time": "2023-06-11T16:34:46.150535-03:00",
  "location": "Jacksonville, FL"
}
```

Running it:

```
make producer FILE_NAME=<path/to/file>
```

### producing a file with random transactions

```
make sample-data TOTAL=1000 FILE_NAME=onethousand.txt
```

## consumer

It listens to a Kafka topic and then process the transaction. If `transaction_amount` is greater than a given value, it is considered as "suspicious" and it is saved to a MongoDB collection.

Running it:

```
make consumer
```

## load testing

To test the consumer with a high number of incoming messages from the topic:

```
make kafka-consumer-publish FILE_NAME=<path/to/file>
```

It uses `kafka-console-producer` script that comes with default instalation of Kafka. It is really fast.

## running tests

```
make test
```

## unit test coverage

```
make coverage
```

## Makefile targets

```
make help
```