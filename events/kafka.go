package events

import (
	"bytes"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/jsonpb"
)

// KafkaConfig describes configuration for accessing a Kafka topic.
type KafkaConfig struct {
	Servers []string
	Topic   string
}

// DefaultKafkaConfig returns default Kafka config.
func DefaultKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Servers: []string{"127.0.0.1:9092"},
		Topic:   "funnel",
	}
}

// KafkaReader allows reading Event messages from a Kafka topic.
type KafkaReader struct {
	con sarama.Consumer
	p   sarama.PartitionConsumer
}

// NewKafkaReader creates a new KafkaReader with the given config.
// TODO currently this is hard-coded to start at the beginning of the topic.
func NewKafkaReader(conf KafkaConfig) (*KafkaReader, error) {
	con, err := sarama.NewConsumer(conf.Servers, nil)
	if err != nil {
		return nil, err
	}

	p, err := con.ConsumePartition(conf.Topic, 0, sarama.OffsetOldest)
	if err != nil {
		return nil, err
	}

	return &KafkaReader{con, p}, nil
}

// Read reads an event from the topic. This blocks until a message is delivered.
func (r *KafkaReader) Read() (*Event, error) {
	msg := <-r.p.Messages()
	ev := &Event{}
	rdr := bytes.NewReader(msg.Value)
	err := jsonpb.Unmarshal(rdr, ev)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal event: %s", err)
	}
	return ev, nil
}
