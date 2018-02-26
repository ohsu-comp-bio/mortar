package cmd

import (
	"fmt"

	"github.com/bmeg/arachne/aql"
	"github.com/ohsu-comp-bio/mortar/graph"
	"github.com/spf13/cobra"
)

func init() {
	conf := DefaultConfig()
	runCmd := cobra.Command{
		Use: "add-step",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("usage: add-step <workflow ID> <step ID>")
			}
			return runAddStep(conf, args[0], args[1])
		},
	}
	cmd.AddCommand(&runCmd)

	f := runCmd.Flags()
	f.StringVar(&conf.Arachne.Server, "Arachne.Server", conf.Arachne.Server, "")
	f.StringVar(&conf.Arachne.Graph, "Arachne.Graph", conf.Arachne.Graph, "")
}

func runAddStep(conf Config, wfid, sid string) error {

	log.Info("Connecting to arachne", "server", conf.Arachne.Server)
	acli, err := aql.Connect(conf.Arachne.Server, true)
	if err != nil {
		return err
	}

	cli := graph.Client{Client: &acli, Graph: conf.Arachne.Graph}
	b := &graph.Batch{}

	step := &graph.Step{ID: sid}
	wf := &graph.Workflow{ID: wfid}

	b.AddVertex(step)
	b.AddVertex(wf)
	b.AddEdge(graph.StepInWorkflow(step, wf))

	err = cli.AddBatch(b)
	if err != nil {
		return err
	}
	return nil
}
