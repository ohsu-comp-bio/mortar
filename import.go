package main

import (
  "github.com/Shopify/sarama"
  "github.com/golang/protobuf/proto"
  "github.com/bmeg/arachne/aql"
  "github.com/ohsu-comp-bio/mortar/events"
  "github.com/ohsu-comp-bio/tes"
  "github.com/ohsu-comp-bio/funnel/config"
  "github.com/ohsu-comp-bio/funnel/logger"
  "github.com/golang/protobuf/jsonpb"
  structpb "github.com/golang/protobuf/ptypes/struct"
)

var log = logger.NewLogger("test", logger.DefaultConfig())
var mar = jsonpb.Marshaler{}

func main() {
  conf := config.Kafka{
    Servers: []string{"127.0.0.1:9092"},
    Topic: "funnel",
  }

	con, err := sarama.NewConsumer(conf.Servers, nil)
	if err != nil {
    panic(err)
	}

	p, err := con.ConsumePartition(conf.Topic, 0, sarama.OffsetOldest)
	if err != nil {
    panic(err)
	}

  cli, err := aql.Connect("127.0.0.1:5757", true)
	if err != nil {
    panic(err)
	}

  graphid := "mortar-13"
  cli.AddGraph(graphid)

  for msg := range p.Messages() {

    ev := &events.Event{}
    err := events.Unmarshal(msg.Value, ev)
    if err != nil {
      log.Error("error unmarshaling event", err)
      continue
    }

    task := &tes.Task{}

    v, err := cli.GetVertex(graphid, ev.Id)
    if err == nil && v != nil {
      unmarshal(v.Data, task)
    }

    events.WriteEvent(task, ev)
    log.Info("task", task)
    st := marshal(task)

    err = cli.AddVertex(graphid, aql.Vertex{
      Gid: ev.Id,
      Label: "Task",
      Data: st,
    })
    if err != nil {
      log.Error("error adding vertex", err)
    }
  }
}

func marshal(msg proto.Message) *structpb.Struct {
  s, err := mar.MarshalToString(msg)
  if err != nil {
    panic(err)
  }

  st := &structpb.Struct{}
  err = jsonpb.UnmarshalString(s, st)
  if err != nil {
    panic(err)
  }

  return st
}

func unmarshal(st *structpb.Struct, msg proto.Message) {
  b, err := mar.MarshalToString(st)
  if err != nil {
    panic(err)
  }

  err = jsonpb.UnmarshalString(b, msg)
  if err != nil {
    panic(err)
  }
}
