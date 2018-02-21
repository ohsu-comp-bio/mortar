package main

import (
	"sort"

	"github.com/bmeg/arachne/aql"
	"github.com/ohsu-comp-bio/mortar/graph"
	"github.com/ohsu-comp-bio/tes"
)

type workflowInfo struct {
	ID    string
	Steps []*graph.Step
}

func getWorkflowInfo(cli *graph.Client, wfID string) *workflowInfo {
	d := &workflowInfo{
		ID:    wfID,
		Steps: []*graph.Step{},
	}

	// Get all steps in the workflow.
	q := aql.V(wfID).In("ktl.StepInWorkflow")
	res, err := cli.Execute(q)
	if err != nil {
		panic(err)
	}

	for row := range res {
		v := row.Value.GetVertex()
		d.Steps = append(d.Steps, &graph.Step{v.Gid})
	}

	return d
}

// Data for a view showing runs (rows) by steps (columns) in a single workflow.
type runsByStep struct {
	Runs []*runStatus
	// The table may be sparse, so this ensures there's a column for every step.
	Steps []*graph.Step
}

func getRunsByStep(cli *graph.Client, wfID string) *runsByStep {
	d := &runsByStep{
		Runs:  []*runStatus{},
		Steps: []*graph.Step{},
	}

	// Get all steps in the workflow.
	q := aql.V(wfID).In("ktl.StepInWorkflow")
	res, err := cli.Execute(q)
	if err != nil {
		panic(err)
	}

	for row := range res {
		v := row.Value.GetVertex()
		d.Steps = append(d.Steps, &graph.Step{v.Gid})
	}

	// Get all the runs for the workflow.
	q = aql.V(wfID).In("ktl.RunForWorkflow")
	res, err = cli.Execute(q)
	if err != nil {
		panic(err)
	}

	for row := range res {
		v := row.Value.GetVertex()
		run := getRunStatus(cli, v.Gid)
		d.Runs = append(d.Runs, run)
	}

	// TODO
	//sort.Sort(d.Steps)
	//sort.Sort(d.Runs)

	return d
}

type runStatus struct {
	ID       string
	Total    int
	Idle     int
	Queued   int
	Running  int
	Error    int
	Complete int
	State    string
	Steps    []*stepStatus
}

func getRunStatus(cli *graph.Client, runID string) *runStatus {
	d := &runStatus{
		ID:    runID,
		Steps: []*stepStatus{},
	}

	// Get all steps.
	stepsQ := aql.V(runID).
		Out("ktl.RunForWorkflow").
		In("ktl.StepInWorkflow")

	res, err := cli.Execute(stepsQ)
	if err != nil {
		panic(err)
	}

	for row := range res {
		v := row.Value.GetVertex()
		s := getStepStatus(cli, runID, v.Gid)
		d.Steps = append(d.Steps, s)

		d.Total++
		switch s.State {
		case tes.Complete:
			d.Complete++
		case tes.Queued:
			d.Queued++
		case tes.Initializing, tes.Running:
			d.Running++
		case tes.ExecutorError, tes.SystemError:
			d.Error++
		}
	}

	d.Idle = d.Total - (d.Complete + d.Queued + d.Running + d.Error)

	switch {
	case d.Error > 0:
		d.State = "error"
	case d.Running > 0:
		d.State = "running"
	case d.Queued > 0:
		d.State = "queued"
	case d.Complete == d.Total:
		d.State = "complete"
	default:
		d.State = "idle"
	}

	return d
}

type stepStatus struct {
	ID     string
	State  tes.State
	Latest *tes.Task
	Tasks  []*tes.Task
}

func getStepStatus(cli *graph.Client, runID, stepID string) *stepStatus {
	d := &stepStatus{
		ID:    stepID,
		Tasks: []*tes.Task{},
	}

	q := aql.V(stepID).
		In("ktl.TaskForStep").As("task").
		Out("ktl.TaskForRun").
		HasID(runID).
		Select("task")

	res, err := cli.Execute(q)
	if err != nil {
		panic(err)
	}

	for row := range res {
		task := &tes.Task{}

		err := graph.Unmarshal(row.Value.GetVertex().Data, task)
		if err != nil {
			panic(err)
		}

		d.Tasks = append(d.Tasks, task)

		// Find the most recent task
		if task.GetCreationTime() > d.Latest.GetCreationTime() {
			d.Latest = task
		}
	}

	// Use the state of the most recent task as the state of the step.
	d.State = d.Latest.GetState()
	return d
}

// Data for a view showing workflows (rows) by a run variable (columns) e.g. "tumor"
// Each cell is a run, showing the summarized status of that run.
type workflowRuns struct {
	Workflows []*workflowStatus
	Columns   []string
}

type workflowStatus struct {
	ID           string
	RunsByColumn map[string]*runStatus
}

func getWorkflowRuns(cli *graph.Client) *workflowRuns {
	// TODO in the future, columnKey will come from the UI
	columnKey := "Sample"

	d := &workflowRuns{
		Workflows: []*workflowStatus{},
		Columns:   []string{},
	}

	// Keep a unique set of column names.
	columns := map[string]bool{}

	// Get all workflows
	q := aql.V().HasLabel("ktl.Workflow")

	res, err := cli.Execute(q)
	if err != nil {
		panic(err)
	}

	for row := range res {
		wfv := row.Value.GetVertex()
		wf := &workflowStatus{
			ID:           wfv.Gid,
			RunsByColumn: map[string]*runStatus{},
		}
		d.Workflows = append(d.Workflows, wf)

		// Get all runs for this workflow
		q := aql.V(wf.ID).In("ktl.RunForWorkflow")

		res, err = cli.Execute(q)
		if err != nil {
			panic(err)
		}

		for row := range res {
			run := row.Value.GetVertex()
			data := map[string]string{}
			graph.Unmarshal(run.Data, &data)
			col, ok := data[columnKey]
			if !ok {
				continue
			}
			columns[col] = true
			wf.RunsByColumn[col] = getRunStatus(cli, run.Gid)
		}
	}

	for col, _ := range columns {
		d.Columns = append(d.Columns, col)
	}

	// TODO
	//sort.Sort(d.Workflows)
	sort.Strings(d.Columns)

	return d
}
