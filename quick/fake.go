package quick

import (
	"github.com/ohsu-comp-bio/mortar/graph"
	"github.com/ohsu-comp-bio/tes"
)

func makeData() *graph.Batch {

	b := &graph.Batch{}

	wf1 := &graph.Workflow{ID: "WF1"}

	s1 := &graph.Step{ID: "S1"} /*, &tes.Task{
		Name:        "S1",
		Description: "Mortar step example",
		Executors: []*tes.Executor{
			{
				Image:   "alpine",
				Command: []string{"echo", "hello", "world"},
			},
		},
	}}
	*/

	s2 := &graph.Step{ID: "S2"} /*, &tes.Task{
		Name:        "S2",
		Description: "Mortar step 2 example",
		Executors: []*tes.Executor{
			{
				Image:   "alpine",
				Command: []string{"md5sum", "example.txt"},
			},
		},
	}}
	*/

	s3 := &graph.Step{ID: "S3"} /*, &tes.Task{
		Name:        "S3",
		Description: "Mortar step 3 example",
		Executors: []*tes.Executor{
			{
				Image:   "alpine",
				Command: []string{"cat"},
			},
		},
	}}
	*/

	r1t1 := &graph.Task{&tes.Task{
		Id:          "run1-123",
		Name:        "Run1 S1",
		Description: "Mortar step example",
		Executors: []*tes.Executor{
			{
				Image:   "alpine",
				Command: []string{"echo", "hello", "world"},
			},
		},
	}}

	r1t2 := &graph.Task{&tes.Task{
		Id:          "run1-124",
		Name:        "Run1 S2",
		Description: "Mortar step 2 example",
		Executors: []*tes.Executor{
			{
				Image:   "alpine",
				Command: []string{"md5sum", "example.txt"},
			},
		},
	}}

	r2t1 := &graph.Task{&tes.Task{
		Id:          "run2-125",
		Name:        "Run2 S1",
		Description: "Mortar step example",
		Executors: []*tes.Executor{
			{
				Image:   "alpine",
				Command: []string{"echo", "hello", "world"},
			},
		},
	}}

	r2t2 := &graph.Task{&tes.Task{
		Id:          "run2-126",
		Name:        "Run2 S2",
		State:       tes.Complete,
		Description: "Mortar step 2 example",
		Executors: []*tes.Executor{
			{
				Image:   "alpine",
				Command: []string{"md5sum", "example.txt"},
			},
		},
	}}

	r2t3 := &graph.Task{&tes.Task{
		Id:          "run2-127",
		Name:        "Run2 S3",
		Description: "Mortar step 3 example",
		Executors: []*tes.Executor{
			{
				Image:   "alpine",
				Command: []string{"md5sum", "example.txt"},
			},
		},
	}}

	r1 := &graph.Run{ID: "Run 1"}
	r2 := &graph.Run{ID: "Run 2"}

	b.AddVertex(wf1)
	b.AddVertex(s1)
	b.AddVertex(s2)
	b.AddVertex(s3)
	b.AddEdge(graph.StepInWorkflow(s1, wf1))
	b.AddEdge(graph.StepInWorkflow(s2, wf1))
	b.AddEdge(graph.StepInWorkflow(s3, wf1))
	b.AddVertex(r1t1)
	b.AddVertex(r1t2)
	b.AddVertex(r2t1)
	b.AddVertex(r2t2)
	b.AddVertex(r2t3)
	b.AddEdge(graph.TaskForStep(r1t1, s1))
	b.AddEdge(graph.TaskForStep(r1t2, s2))
	b.AddEdge(graph.TaskForStep(r2t1, s1))
	b.AddEdge(graph.TaskForStep(r2t2, s2))
	b.AddEdge(graph.TaskForStep(r2t3, s3))
	b.AddVertex(r1)
	b.AddVertex(r2)
	b.AddEdge(graph.RunForWorkflow(r1, wf1))
	b.AddEdge(graph.RunForWorkflow(r2, wf1))
	b.AddEdge(graph.TaskForRun(r1t1, r1))
	b.AddEdge(graph.TaskForRun(r1t2, r1))
	b.AddEdge(graph.TaskForRun(r2t1, r2))
	b.AddEdge(graph.TaskForRun(r2t2, r2))
	b.AddEdge(graph.TaskForRun(r2t3, r2))

	return b
}
