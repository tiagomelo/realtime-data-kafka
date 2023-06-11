SHELL = /bin/bash

include .env
export

SAMPLE_DATA_FOLDER=sampledata

# ==============================================================================
# Help

.PHONY: help
## help: shows this help message
help:
	@ echo "Usage: make [target]"
	@ sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# ==============================================================================
# Kafka

.PHONY: zookeeper
## zookeeper: starts zookeeper
zookeeper: 
	@ zookeeper-server-start /opt/homebrew/etc/zookeeper/zoo.cfg

.PHONY: kafka
## kafka: starts kafka
kafka:
	@ kafka-server-start /opt/homebrew/etc/kafka/server.properties

.PHONY: kafka-consumer-publish
## kafka-consumer-publish: Kafka's tool to read data from standard input and publish it to Kafka
kafka-consumer-publish:
	@ if [ -z "$(FILE_NAME)" ]; then echo >&2 please set file name via the variable FILE_NAME; exit 2; fi
	@ cat $(FILE_NAME) | kafka-console-producer --topic $(KAFKA_TOPIC) --bootstrap-server $(KAFKA_BROKER_HOST)

.PHONY: clear-kafka-messages
## clear-kafka-messages: cleans all pending messages from Kafka
clear-kafka-messages:
	@ rm -rf /opt/homebrew/var/lib/kafka-logs/*

# ==============================================================================
# Producer

.PHONY: producer
## producer: starts producer
producer:
	@ if [ -z "$(FILE_NAME)" ]; then echo >&2 please set file name via the variable FILE_NAME; exit 2; fi
	@ go run producer/producer.go -f=$(FILE_NAME)

# ==============================================================================
# Consumer

.PHONY: consumer
## consumer: starts consumer
consumer:
	@ go run consumer/consumer.go

# ==============================================================================
# Tests

.PHONY: test
## test: runs tests
test:
	@ go test -v ./...

.PHONY: coverage
## coverage: run unit tests and generate coverage report in html format
coverage:
	@ go test -coverprofile=coverage.out ./...  && go tool cover -html=coverage.out

# ==============================================================================
# Sample data generation

.PHONY: sample-data
## sample-data: generates sample data
sample-data:
	@ if [ -z "$(TOTAL)" ]; then echo >&2 please set total via the variable TOTAL; exit 2; fi
	@ if [ -z "$(FILE_NAME)" ]; then echo >&2 please set file name via the variable FILE_NAME; exit 2; fi
	@ rm -f "${SAMPLE_DATA_FOLDER}/${FILE_NAME}"
	@ echo "generating file ${SAMPLE_DATA_FOLDER}/${FILE_NAME}..."
	@ go run jsongenerator/jsongenerator.go --llmin 10000 --llmax 30000 --ulmin 100 --ulmax 3000 -t=$(TOTAL) -p=0.7 -f="${SAMPLE_DATA_FOLDER}/${FILE_NAME}"
	@ echo "file ${SAMPLE_DATA_FOLDER}/${FILE_NAME} was generated." 