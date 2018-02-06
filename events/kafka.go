package events

import (
  "fmt"
  "github.com/Shopify/sarama"
)

type KafkaConfig struct {
  Servers []string
  Topic string
}

func DefaultKafkaConfig() KafkaConfig {
  return KafkaConfig{
    Servers: []string{"127.0.0.1:9092"},
    Topic: "funnel",
  }
}

type KafkaReader struct {
  con sarama.Consumer
  p sarama.PartitionConsumer
}

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

func (r *KafkaReader) Read() (*Event, error) {
  msg := <-r.p.Messages()
  ev := &Event{}
  err := Unmarshal(msg.Value, ev)
  if err != nil {
    return nil, fmt.Errorf("can't unmarshal event: %s", err)
  }
  return ev, nil
}
