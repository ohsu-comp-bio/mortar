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
		Use:  "run",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(conf)
		},
	}
	cmd.AddCommand(&runCmd)

	f := runCmd.Flags()
	f.StringVar(&conf.Arachne.Server, "Arachne.Server", conf.Arachne.Server, "")
	f.StringVar(&conf.Arachne.Graph, "Arachne.Graph", conf.Arachne.Graph, "")
	f.StringSliceVar(&conf.Kafka.Servers, "Kafka.Servers", conf.Kafka.Servers, "")
	f.StringVar(&conf.Kafka.Topic, "Kafka.Topic", conf.Kafka.Topic, "")
}

func run(conf Config) error {

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
	counter := NewCounter("imported events", time.Second)

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

		task := &tes.Task{Id: ev.Id}

		v, err := cli.GetVertex(ev.Id)
		if err == nil && v != nil && v.Data != nil {
			graph.Unmarshal(v.Data, task)
		}

		events.WriteEvent(task, ev)

		taskV := &graph.Task{Task: task}
		err = cli.AddVertex(taskV)
		if err != nil {
			log.Error("error adding vertex", err)
			continue
		}

		switch ev.Type {
		case events.Type_TASK_CREATED:

			if sID, ok := task.Tags["ktl.StepID"]; ok {
				step := &graph.Step{ID: sID}
				err = cli.AddEdge(graph.TaskForStep(taskV, step))
				if err != nil {
					log.Error("error adding edge", err)
				}
			}

			for _, input := range task.Inputs {
				// Some inputs have an empty URL, e.g. if they define the "content" field
				if input.Url == "" {
					continue
				}

				iv := &graph.File{input.Url, input.Type}
				err := cli.AddVertex(iv)
				if err != nil {
					log.Error("error adding vertex", err)
				}

				err = cli.AddEdge(graph.TaskRequestsInput(taskV, iv))
				if err != nil {
					log.Error("can't add edge", err)
				}
			}

			for _, output := range task.Outputs {
				ov := &graph.File{output.Url, output.Type}
				err := cli.AddVertex(ov)
				if err != nil {
					log.Error("error adding vertex", err)
				}

				err = cli.AddEdge(graph.TaskRequestsOutput(taskV, ov))
				if err != nil {
					log.Error("can't add edge", err)
				}
			}

			for _, exec := range task.Executors {
				iv := &graph.Image{exec.Image}
				err := cli.AddVertex(iv)
				if err != nil {
					log.Error("error adding vertex", err)
				}

				err = cli.AddEdge(graph.TaskRequestsImage(taskV, iv))
				if err != nil {
					log.Error("can't add edge", err)
				}
			}

		case events.Type_TASK_OUTPUTS:
			outputs := ev.GetOutputs().Value
			for _, output := range outputs {
				ov := &graph.File{URL: output.Url}
				err := cli.AddVertex(ov)
				if err != nil {
					log.Error("error adding vertex", err)
				}

				err = cli.AddEdge(graph.TaskUploadedOutput(taskV, ov))
				if err != nil {
					log.Error("can't add edge", err)
				}
			}
		}
		counter.Inc()
	}
}
