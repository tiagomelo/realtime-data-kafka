package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

// Config holds all configuration needed by this app.
type Config struct {
	KafkaBrokerHost string `envconfig:"KAFKA_BROKER_HOST" required:"true"`
	KafkaTopic      string `envconfig:"KAFKA_TOPIC" required:"true"`
	KafkaGroupId    string `envconfig:"KAFKA_GROUP_ID" required:"true"`
	MongodbDatabase string `envconfig:"MONGODB_DATABASE" required:"true"`
	MongodbHostName string `envconfig:"MONGODB_HOST_NAME" required:"true"`
	MongodbPort     int    `envconfig:"MONGODB_PORT" required:"true"`
}

// For ease of unit testing.
var (
	godotenvLoad     = godotenv.Load
	envconfigProcess = envconfig.Process
)

// Read reads the environment variables from the given file and returns a Config.
func Read(envFile string) (*Config, error) {
	if err := godotenvLoad(envFile); err != nil {
		return nil, errors.Wrap(err, "loading env vars")
	}
	config := new(Config)
	if err := envconfigProcess("", config); err != nil {
		return nil, errors.Wrap(err, "processing env vars")
	}
	return config, nil
}
