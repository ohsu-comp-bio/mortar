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
	runCmd := cobra.Command{
		Use:  "runsmc",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSmchet(conf)
		},
	}
	cmd.AddCommand(&runCmd)

	f := runCmd.Flags()
	f.StringVar(&conf.Arachne.Server, "Arachne.Server", conf.Arachne.Server, "")
	f.StringVar(&conf.Arachne.Graph, "Arachne.Graph", conf.Arachne.Graph, "")
	f.StringSliceVar(&conf.Kafka.Servers, "Kafka.Servers", conf.Kafka.Servers, "")
	f.StringVar(&conf.Kafka.Topic, "Kafka.Topic", conf.Kafka.Topic, "")
}

func runSmchet(conf Config) error {

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
	counter := NewCounter("imported smc-het events", time.Second)

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

		task := &graph.Task{Task: &tes.Task{Id: ev.Id}}

		v, err := cli.GetVertex(ev.Id)
		if err == nil && v != nil && v.Data != nil {
			graph.Unmarshal(v.Data, task.Task)
		}

		events.WriteEvent(task.Task, ev)

		// Skip events that aren't creating a task.
		if ev.Type != events.Type_TASK_CREATED {
			continue
		}

		enID := task.Tags["entry"]
		evID := task.Tags["eval"]
		tuID := task.Tags["tumor"]
		tnID := task.Tags["taskname"]

		// Try to make up for missing "eval" tag in legacy data
		if evID == "" && enID != "" && tuID != "" {
			evID = enID + "/" + tuID
		}

		// If missing required data, skip.
		if tnID == "" || evID == "" || enID == "" {
			continue
		}

		run := &graph.Run{ID: evID}
		wf := &graph.Workflow{ID: enID}
		step := &graph.Step{ID: tnID}

		b := &graph.Batch{}
		b.AddVertex(run)
		b.AddVertex(step)
		b.AddVertex(wf)
		b.AddEdge(graph.TaskForStep(task, step))
		b.AddEdge(graph.TaskForRun(task, run))
		b.AddEdge(graph.StepInWorkflow(step, wf))
		b.AddEdge(graph.RunForWorkflow(run, wf))

		err = cli.AddBatch(b)
		if err != nil {
			log.Error("add batch failed", err)
			continue
		}

		counter.Inc()
	}
}
