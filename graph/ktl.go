package graph

import (
	"fmt"

	"github.com/bmeg/arachne/aql"
	"github.com/ohsu-comp-bio/tes"
)

// TODO need names for the difference between workflow descriptor and invocation

type Run struct {
  ID string
  // TODO want state string? so that it's easier to check the status of only
  //      running runs?
}

func (r *Run) MarshalAQL() (*aql.Vertex, error) {
  return &aql.Vertex{
    Gid: r.ID,
    Label: "ktl.Run",
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
}

func (w *Workflow) MarshalAQL() (*aql.Vertex, error) {
	return &aql.Vertex{
		Gid:   w.ID,
		Label: "ktl.Workflow",
	}, nil
}

func StepInWorkflow(step *Step, wf *Workflow) Edge {
	return NewEdge("ktl.StepInWorkflow", step, wf)
}

// Step describes a step in a ktl workflow.
type Step struct {
	ID string
  // Reusing tes.Task to describe a step,
  // but stateful fields (state, logs, execution data)
  // should be ignored.
  // Possibly just a placeholder?
  // Does TES need a stateless task description?
	*tes.Task
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

func TaskForStep(task *Task, step *Step) Edge {
	return NewEdge("ktl.TaskForStep", task, step)
}
