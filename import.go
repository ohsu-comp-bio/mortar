package main

import (
  "context"
  "github.com/Shopify/sarama"
  "github.com/bmeg/arachne/aql"
  "github.com/ohsu-comp-bio/mortar/events"
  "github.com/ohsu-comp-bio/tes"
  "github.com/ohsu-comp-bio/funnel/config"
  "github.com/ohsu-comp-bio/funnel/logger"
  "github.com/golang/protobuf/jsonpb"
  structpb "github.com/golang/protobuf/ptypes/struct"
  "time"
)

func main() {
  conf := config.Kafka{
    Servers: []string{"10.50.50.85:9092", "10.50.50.84:9092", "10.50.50.83:9092"},
    Topic: "funnel-events",
  }
  log := logger.NewLogger("test", logger.DefaultConfig())

	con, err := sarama.NewConsumer(conf.Servers, nil)
	if err != nil {
    panic(err)
	}

	p, err := con.ConsumePartition(conf.Topic, 0, sarama.OffsetOldest)
	if err != nil {
    panic(err)
	}

  cli, err := aql.Connect("10.50.50.123:9090", true)
	if err != nil {
    panic(err)
	}
  log.Info("graphs", cli.GetGraphList())

  ctx := context.Background()
  _ = ctx

  builder := events.TaskBuilder{}
  mar := jsonpb.Marshaler{}

  graphid := "mortar-6"
  cli.AddGraph(graphid)

  for msg := range p.Messages() {
    ev := &events.Event{}
    err := events.Unmarshal(msg.Value, ev)
    if err != nil {
      log.Error("error unmarshaling event", err)
      continue
    }

    v, _ := cli.GetVertex(graphid, ev.Id)
    log.Info("getting vertex", "id", ev.Id, "vert", v)
    /*
    if err != nil {
      log.Error("error getting vertex", err)
      continue
    }
    */

    task := &tes.Task{}
    if v != nil {
      s, err := mar.MarshalToString(v.Data)
      if err != nil {
        log.Error("error marshaling data")
        continue
      }
      err = jsonpb.UnmarshalString(s, task)
      if err != nil {
        log.Error("error unmarshaling task from string")
        continue
      }
      log.Info("loaded vertex", "task", task)
    }

    builder.Task = task
    err = builder.WriteEvent(ctx, ev)
    if err != nil {
      log.Error("error building task")
      continue
    }

    b, err := mar.MarshalToString(task)
    if err != nil {
      log.Error("error marshaling task")
      continue
    }

    s := &structpb.Struct{}
    err = jsonpb.UnmarshalString(b, s)
    if err != nil {
      log.Error("error unmarshaling to struct type")
      continue
    }

    log.Info("event", ev)
    log.Info("task", s)
    err = cli.AddVertex(graphid, aql.Vertex{
      Gid: ev.Id,
      Label: "Task",
      Data: s,
    })
    if err != nil {
      log.Error("error adding vertex", err)
    }
    time.Sleep(100 * time.Millisecond)
  }
}
