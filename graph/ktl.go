package graph

import (
	"fmt"

	"github.com/bmeg/arachne/aql"
	"github.com/ohsu-comp-bio/tes"
)

type Workflow struct {
	ID string
}

func (w *Workflow) MarshalAQL() (*aql.Vertex, error) {
	return &aql.Vertex{
		Gid:   w.ID,
		Label: "ktl.Workflow",
	}, nil
}

func StepInWorkflow(step *Step, wf *Workflow) Edge {
	return NewEdge("kt.StepInWorkflow", step, wf)
}

// Step describes a step in a ktl workflow.
// It happens to be the same as a tes.Task, for now,
// without any state/logs/execution data.
type Step struct {
	ID string
	*tes.Task
}

// MarshalAQL marshals the vertex into an arachne AQL vertex.
func (s *Step) MarshalAQL() (*aql.Vertex, error) {
	d, err := Marshal(s.Task)
	if err != nil {
		return nil, fmt.Errorf("can't marshal ktl.Step: %s", err)
	}
	return &aql.Vertex{
		Gid:   s.ID,
		Label: "ktl.Step",
		Data:  d,
	}, nil
}

func TaskForStep(task *Task, step *Step) Edge {
	return NewEdge("ktl.TaskForStep", task, step)
}
