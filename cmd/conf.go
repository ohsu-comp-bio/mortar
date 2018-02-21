package cmd

import (
	"github.com/ohsu-comp-bio/mortar/events"
)

// Config describes the top-level configuration.
type Config struct {
	Kafka   events.KafkaConfig
	Arachne ArachneConfig
}

// DefaultConfig returns the default config.
func DefaultConfig() Config {
	return Config{
		Kafka:   events.DefaultKafkaConfig(),
		Arachne: DefaultArachneConfig(),
	}
}

// ArachneConfig describes the arachne database configuration.
type ArachneConfig struct {
	Server string
	Graph  string
}

// DefaultArachneConfig returns the default arachne database configuration.
func DefaultArachneConfig() ArachneConfig {
	return ArachneConfig{
		Server: "localhost:8202",
		Graph:  "mortar",
	}
}
