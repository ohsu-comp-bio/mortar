package cmd

import (
	"time"

	"github.com/bmeg/arachne/aql"
	"github.com/ohsu-comp-bio/mortar/events"
	"github.com/ohsu-comp-bio/mortar/graph"
	"github.com/ohsu-comp-bio/tes"
	"github.com/spf13/cobra"
)

func init() {
	conf := DefaultConfig()
	icmd := &cobra.Command{
		Use:  "ingest",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return ingest(conf)
		},
	}
	cmd.AddCommand(icmd)

	f := icmd.Flags()
	f.StringVar(&conf.Arachne.Server, "Arachne.Server", conf.Arachne.Server, "")
	f.StringVar(&conf.Arachne.Graph, "Arachne.Graph", conf.Arachne.Graph, "")
	f.StringSliceVar(&conf.Kafka.Servers, "Kafka.Servers", conf.Kafka.Servers, "")
	f.StringVar(&conf.Kafka.Topic, "Kafka.Topic", conf.Kafka.Topic, "")
}

func ingest(conf Config) error {

	log.Info("Connecting to arachne", "server", conf.Arachne.Server)
	acli, err := aql.Connect(conf.Arachne.Server, true)
	if err != nil {
		return err
	}

	err = acli.AddGraph(conf.Arachne.Graph)
	if err != nil {
		return err
	}

	cli := graph.Client{Client: &acli, Graph: conf.Arachne.Graph}
	counter := NewCounter("ingested events", time.Second)
	readCounter := NewCounter("read events", time.Second)

	r, err := events.NewKafkaReader(conf.Kafka)
	if err != nil {
		return err
	}

	for {
		ev, err := r.Read()
		if err != nil {
			log.Error("can't read event", err)
			continue
		}
		readCounter.Inc()

		task := &graph.Task{Task: &tes.Task{Id: ev.Id}}

		v, err := cli.GetVertex(ev.Id)
		if err == nil && v != nil && v.Data != nil {
			graph.Unmarshal(v.Data, task.Task)
		}

		events.WriteEvent(task.Task, ev)
		b := &graph.Batch{}
		b.AddVertex(task)

		switch ev.Type {
		case events.Type_TASK_CREATED:

			// TODO step IDs need to be globally unique
			stepID := task.Tags["mortar.StepID"]
			runID := task.Tags["mortar.RunID"]

			if stepID == "" {
				log.Error("missing stepID")
				continue
			}
			if runID == "" {
				log.Error("missing runID")
				continue
			}

			// TODO should these create vertices?

			step := &graph.Step{ID: stepID}
			run := &graph.Run{ID: runID}
			b.AddEdge(graph.TaskForStep(task, step))
			b.AddEdge(graph.TaskForRun(task, run))
		}

		err = cli.AddBatch(b)
		if err != nil {
			log.Error("add batch failed", err)
			continue
		}

		counter.Inc()
	}
}
