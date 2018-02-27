package graph

import (
	"fmt"

	"github.com/bmeg/arachne/aql"
	//"github.com/ohsu-comp-bio/tes"
)

// TODO need names for the difference between workflow descriptor and invocation

type Run struct {
	ID     string
  Inputs map[string]interface{}
}

func (r *Run) MarshalAQL() (*aql.Vertex, error) {
	d, _ := Marshal(r)

	return &aql.Vertex{
		Gid:   r.ID,
		Label: "ktl.Run",
		Data:  d,
	}, nil
}

func TaskForRun(task *Task, run *Run) Edge {
	return NewEdge("ktl.TaskForRun", task, run)
}

func RunForWorkflow(run *Run, wf *Workflow) Edge {
	return NewEdge("ktl.RunForWorkflow", run, wf)
}

type Workflow struct {
	ID string
  Doc map[string]interface{}
}

func (w *Workflow) MarshalAQL() (*aql.Vertex, error) {
	d, _ := Marshal(w)

	return &aql.Vertex{
		Gid:   w.ID,
		Label: "ktl.Workflow",
    Data: d,
	}, nil
}

func StepInWorkflow(step *Step, wf *Workflow) Edge {
	return NewEdge("ktl.StepInWorkflow", step, wf)
}

// Step describes a step in a ktl workflow.
type Step struct {
	ID string
  Doc map[string]interface{}
  Order int
}

// MarshalAQL marshals the vertex into an arachne AQL vertex.
func (s *Step) MarshalAQL() (*aql.Vertex, error) {
	d, err := Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("can't marshal ktl.Step: %s", err)
	}
	return &aql.Vertex{
		Gid:   s.ID,
		Label: "ktl.Step",
		Data:  d,
	}, nil
}

func (s *Step) UnmarshalAQL(v *aql.Vertex) error {
  return Unmarshal(v.Data, s)
}

type OrderedSteps []*Step
func (o OrderedSteps) Len() int { return len(o) }
func (o OrderedSteps) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
func (o OrderedSteps) Less(i, j int) bool { return o[i].Order < o[j].Order }

func TaskForStep(task *Task, step *Step) Edge {
	return NewEdge("ktl.TaskForStep", task, step)
}
