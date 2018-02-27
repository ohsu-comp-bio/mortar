package cmd

import (
	"fmt"
	"crypto/md5"
  "encoding/json"
  "io/ioutil"
  "os"

	"github.com/bmeg/arachne/aql"
	"github.com/ohsu-comp-bio/mortar/graph"
	"github.com/spf13/cobra"
)

func init() {
	conf := DefaultConfig()
	runCmd := cobra.Command{
		Use: "add-wf",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("usage: add-wf <workflow ID> <filepath>")
			}
			return runAddWf(conf, args[0], args[1])
		},
	}
	cmd.AddCommand(&runCmd)

	f := runCmd.Flags()
	f.StringVar(&conf.Arachne.Server, "Arachne.Server", conf.Arachne.Server, "")
	f.StringVar(&conf.Arachne.Graph, "Arachne.Graph", conf.Arachne.Graph, "")
}

func runAddWf(conf Config, wfid, path string) error {
  fh, err := os.Open(path)
  if err != nil {
    return err
  }

  b, err := ioutil.ReadAll(fh)
  if err != nil {
    return err
  }
  log.Info("wf", string(b))

  // Load bunny-style resolved workflow from json
  doc := map[string]interface{}{}
  err = json.Unmarshal(b, &doc)
  if err != nil {
    return err
  }

  bat := &graph.Batch{}

	wf := &graph.Workflow{ID: wfid, Doc: doc}
  bat.AddVertex(wf)

  steps := doc["steps"].([]interface{})
  for i, stepi := range steps {
    step := stepi.(map[string]interface{})
    id := step["id"].(string)
    s := &graph.Step{
      ID: id,
      Doc: step,
      Order: i,
    }
    log.Info("step", "id", id, "step", step, "stepv", s)
    bat.AddVertex(s)
    bat.AddEdge(graph.StepInWorkflow(s, wf))
  }

	log.Info("Connecting to arachne", "server", conf.Arachne.Server)
	acli, err := aql.Connect(conf.Arachne.Server, true)
	if err != nil {
		return err
	}

	cli := graph.Client{Client: &acli, Graph: conf.Arachne.Graph}

  return cli.AddBatch(bat)
}

// hashDoc returns the md5 hexadecimal checksum of the given string.
// used to create a content-based ID of a workflow document.
func hashDoc(s string) string {
  return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
