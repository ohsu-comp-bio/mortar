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
		Use: "add-run",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return fmt.Errorf("usage: add-run <workflow ID> <run ID> <sample>")
			}
			return runAddRun(conf, args[0], args[1], args[2])
		},
	}
	cmd.AddCommand(&runCmd)

	f := runCmd.Flags()
	f.StringVar(&conf.Arachne.Server, "Arachne.Server", conf.Arachne.Server, "")
	f.StringVar(&conf.Arachne.Graph, "Arachne.Graph", conf.Arachne.Graph, "")
}

func runAddRun(conf Config, wfid, rid, sample string) error {

	log.Info("Connecting to arachne", "server", conf.Arachne.Server)
	acli, err := aql.Connect(conf.Arachne.Server, true)
	if err != nil {
		return err
	}

	cli := graph.Client{Client: &acli, Graph: conf.Arachne.Graph}
	b := &graph.Batch{}

	run := &graph.Run{ID: rid, Sample: sample}
	wf := &graph.Workflow{ID: wfid}

	b.AddVertex(run)
	b.AddEdge(graph.RunForWorkflow(run, wf))

	err = cli.AddBatch(b)
	if err != nil {
		return err
	}
	return nil
}
