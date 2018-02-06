package main

import (
  "github.com/bmeg/arachne/aql"
  "github.com/ohsu-comp-bio/mortar/events"
  "github.com/ohsu-comp-bio/mortar/graph"
  "github.com/ohsu-comp-bio/tes"
  "github.com/ohsu-comp-bio/funnel/logger"
)

var log = logger.NewLogger("test", logger.DefaultConfig())

func main() {

  r, err := events.NewKafkaReader(events.DefaultKafkaConfig())
	if err != nil {
    panic(err)
	}

  acli, err := aql.Connect("127.0.0.1:5757", true)
	if err != nil {
    panic(err)
	}
  graphid := "mortar-15"
  acli.AddGraph(graphid)
  cli := graph.Client{Client: &acli, Graph: graphid}

  for {
    ev, err := r.Read()
    if err != nil {
      log.Error("can't read event", err)
      continue
    }

    task := &tes.Task{}

    v, err := cli.GetVertex(ev.Id)
    if err == nil && v != nil && v.Data != nil {
      log.Info("unmarshal data", v.Data)
      graph.Unmarshal(v.Data, task)
    }

    events.WriteEvent(task, ev)
    log.Info("task", task)

    taskV := &graph.TaskVertex{Task: task}
    err = cli.AddVertex(taskV)
    if err != nil {
      log.Error("error adding vertex", err)
      continue
    }

    switch ev.Type {
    case events.Type_TASK_CREATED:

      for k, v := range task.Tags {
        tv := &graph.TagVertex{k, v}
        err := cli.AddVertex(tv)
        if err != nil {
          log.Error("error adding vertex", err)
        }

        err = cli.AddEdge(&graph.HasTagEdge{taskV, tv})
        if err != nil {
          log.Error("error adding task->tag edge", err)
        }
      }

      for _, input := range task.Inputs {
        iv := &graph.FileVertex{input.Url, input.Type}
        err := cli.AddVertex(iv)
        if err != nil {
          log.Error("error adding vertex", err)
        }

        err = cli.AddEdge(&graph.RequestsInputEdge{taskV, iv})
        if err != nil {
          log.Error("can't add edge", err)
        }
      }

      for _, output := range task.Outputs {
        ov := &graph.FileVertex{output.Url, output.Type}
        err := cli.AddVertex(ov)
        if err != nil {
          log.Error("error adding vertex", err)
        }

        err = cli.AddEdge(&graph.RequestsOutputEdge{taskV, ov})
        if err != nil {
          log.Error("can't add edge", err)
        }
      }

      for _, exec := range task.Executors {
        iv := &graph.ImageVertex{exec.Image}
        err := cli.AddVertex(iv)
        if err != nil {
          log.Error("error adding vertex", err)
        }

        err = cli.AddEdge(&graph.RequestsImageEdge{taskV, iv})
        if err != nil {
          log.Error("can't add edge", err)
        }
      }

    case events.Type_TASK_OUTPUTS:
      outputs := ev.GetOutputs().Value
      for _, output := range outputs {
        ov := &graph.FileVertex{Url: output.Url}
        err := cli.AddVertex(ov)
        if err != nil {
          log.Error("error adding vertex", err)
        }

        err = cli.AddEdge(&graph.UploadedOutputEdge{taskV, ov})
        if err != nil {
          log.Error("can't add edge", err)
        }
      }
    }
  }
}
