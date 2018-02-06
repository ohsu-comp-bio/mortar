package cmd

import (
  "github.com/ohsu-comp-bio/mortar/events"
)

type Config struct {
  Kafka events.KafkaConfig
  Arachne ArachneConfig
}

func DefaultConfig() Config {
  return Config{
    Kafka: events.DefaultKafkaConfig(),
    Arachne: DefaultArachneConfig(),
  }
}

type ArachneConfig struct {
  Server string
  Graph string
}

func DefaultArachneConfig() ArachneConfig {
  return ArachneConfig{
    Server: "localhost:5757",
    Graph: "mortar",
  }
}
