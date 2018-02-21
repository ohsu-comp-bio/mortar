package events

import (
	"bytes"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/jsonpb"
)

var mar = jsonpb.Marshaler{}

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

// KafkaWriter writes events to a Kafka topic.
type KafkaWriter struct {
	conf     KafkaConfig
	producer sarama.SyncProducer
}

// NewKafkaWriter creates a new event writer for writing events to a Kafka topic.
func NewKafkaWriter(conf KafkaConfig) (*KafkaWriter, error) {
	producer, err := sarama.NewSyncProducer(conf.Servers, nil)
	if err != nil {
		return nil, err
	}
	return &KafkaWriter{conf, producer}, nil
}

// WriteEvent writes the event. Events may be sent in batches in the background by the
// Kafka client library. Currently stdout, stderr, and system log events are dropped.
func (k *KafkaWriter) Write(ev *Event) error {

	switch ev.Type {
	case Type_EXECUTOR_STDOUT, Type_EXECUTOR_STDERR, Type_SYSTEM_LOG:
		return nil
	}

	s, err := mar.MarshalToString(ev)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: k.conf.Topic,
		Key:   nil,
		Value: sarama.StringEncoder(s),
	}
	_, _, err = k.producer.SendMessage(msg)
	return err
}
