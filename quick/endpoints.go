package quick

import (
	"sort"

	"github.com/bmeg/arachne/aql"
	"github.com/bmeg/arachne/protoutil"
	"github.com/ohsu-comp-bio/mortar/graph"
	"github.com/ohsu-comp-bio/tes"
)

type workflowInfo struct {
	ID    string
  Workflow *graph.Workflow
	Steps []*graph.Step
}

func getWorkflowInfo(cli *graph.Client, wfID string) *workflowInfo {
	d := &workflowInfo{
		ID:    wfID,
		Steps: []*graph.Step{},
	}

  wfv, err := cli.GetVertex(wfID)
  if err != nil {
    log.Error("error", err)
    return nil
  }

  data := protoutil.AsMap(wfv.Data)
  doc, ok := data["Doc"].(map[string]interface{})
  if !ok {
    log.Error("error", err)
    return nil
  }

  d.Workflow = &graph.Workflow{
    ID: wfID,
    Doc: doc,
  }

	// Get all steps in the workflow.
	q := aql.V(wfID).In("ktl.StepInWorkflow")
	res, err := cli.Execute(q)
	if err != nil {
		panic(err)
	}

	for row := range res {
		v := row.Value.GetVertex()
		d.Steps = append(d.Steps, &graph.Step{ID: v.Gid})
	}

	return d
}

// Data for a view showing workflows (rows) by a run variable (columns) e.g. "tumor"
// Each cell is a run, showing the summarized status of that run.
type workflowRuns struct {
	Rows    []string
	Columns []string
	// cell key is "{workflow.ID}-{column}"
	Cells map[string]*runStatus
}

func getWorkflowRuns(cli *graph.Client) *workflowRuns {

	columnKey := "Sample"
	// Track unique column values
	columns := map[string]bool{}

	d := &workflowRuns{
		Rows:    []string{},
		Columns: []string{},
		Cells:   map[string]*runStatus{},
	}

	st := getWorkflowStatuses(cli)

	for _, wf := range st {
		d.Rows = append(d.Rows, wf.ID)

		for _, run := range wf.Runs {
			col, ok := run.Data[columnKey]
			if !ok {
				continue
			}
			columns[col] = true

			cellKey := wf.ID + "-" + col
			d.Cells[cellKey] = run
			run.Steps = nil
		}
	}

	for col, _ := range columns {
		d.Columns = append(d.Columns, col)
	}
	sort.Strings(d.Rows)
	sort.Strings(d.Columns)

	return d
}

type workflowStatus struct {
	ID    string
	Runs  map[string]*runStatus
	Steps []*graph.Step
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
	Steps    map[string]*stepStatus
	Data     map[string]string
}

type stepStatus struct {
	ID     string
	State  tes.State
	Latest *tes.Task
	Tasks  []*tes.Task
}

func getWorkflowStatuses(cli *graph.Client) map[string]*workflowStatus {
	d := map[string]*workflowStatus{}

	// Get all workflows
	q := aql.V().
		HasLabel("ktl.Workflow").As("wf").
		In("ktl.StepInWorkflow").As("step").
		Select("wf", "step")

	res, err := cli.Execute(q)
	if err != nil {
		panic(err)
	}

	for row := range res {
		wfv := row.Row[0].GetVertex()
		stepv := row.Row[1].GetVertex()
    step := &graph.Step{}
    step.UnmarshalAQL(stepv)

		wfst, ok := d[wfv.Gid]
		if !ok {
			wfst = &workflowStatus{
				ID:   wfv.Gid,
				Runs: map[string]*runStatus{},
			}
			d[wfst.ID] = wfst
		}

		wfst.Steps = append(wfst.Steps, step)
	}

  for _, wfst := range d {
    sort.Sort(graph.OrderedSteps(wfst.Steps))
  }

	// Get all runs
	q = aql.V().
		HasLabel("ktl.Workflow").As("wf").
		In("ktl.RunForWorkflow").As("run").
		Select("wf", "run")

	res, err = cli.Execute(q)
	if err != nil {
		panic(err)
	}

	for row := range res {
		wfv := row.Row[0].GetVertex()
		runv := row.Row[1].GetVertex()
		wfst := d[wfv.Gid]
		runst := &runStatus{
			ID:    runv.Gid,
			Steps: map[string]*stepStatus{},
			Total: len(wfst.Steps),
			Data:  map[string]string{},
		}

		graph.Unmarshal(runv.Data, &runst.Data)

		wfst.Runs[runst.ID] = runst
		for _, step := range wfst.Steps {
			runst.Steps[step.ID] = &stepStatus{
				ID:    step.ID,
				Tasks: []*tes.Task{},
			}
		}
	}

	// Get the state of all runs
	q = aql.V().
		HasLabel("ktl.Workflow").As("wf").
		In("ktl.RunForWorkflow").As("run").
		In("ktl.TaskForRun").As("task").
		Out("ktl.TaskForStep").As("step").
		Select("wf", "run", "task", "step")

	res, err = cli.Execute(q)
	if err != nil {
		panic(err)
	}

	for row := range res {
		wfv := row.Row[0].GetVertex()
		runv := row.Row[1].GetVertex()
		taskv := row.Row[2].GetVertex()
		stepv := row.Row[3].GetVertex()

		wfst := d[wfv.Gid]
		runst := wfst.Runs[runv.Gid]
		stepst := runst.Steps[stepv.Gid]

		task := &tes.Task{}

		err := graph.Unmarshal(taskv.Data, task)
		if err != nil {
			panic(err)
		}

		stepst.Tasks = append(stepst.Tasks, task)
	}

	for _, wf := range d {
		for _, run := range wf.Runs {
			for _, step := range run.Steps {
				step.Latest = LatestTask(step.Tasks)
				step.State = step.Latest.GetState()

				// Tally the counts of each step's state.
				switch step.State {
				case tes.Complete:
					run.Complete++
				case tes.Queued:
					run.Queued++
				case tes.Initializing, tes.Running:
					run.Running++
				case tes.ExecutorError, tes.SystemError:
					run.Error++
				}
			}

			run.Idle = run.Total - (run.Complete + run.Queued + run.Running + run.Error)

			// Determine the run state based on the step state counts.
			switch {
			case run.Error > 0:
				run.State = "error"
			case run.Running > 0:
				run.State = "running"
			case run.Queued > 0:
				run.State = "queued"
			case run.Complete == run.Total:
				run.State = "complete"
			default:
				run.State = "idle"
			}
		}
	}

	return d
}

func LatestTask(tasks []*tes.Task) *tes.Task {
	if len(tasks) == 0 {
		return nil
	}
	if len(tasks) == 1 {
		return tasks[0]
	}
	l := make([]*tes.Task, len(tasks))
	copy(l, tasks)
	sort.Sort(ByTaskCreationTime(l))
	return l[len(l)-1]
}

type ByTaskCreationTime []*tes.Task

func (b ByTaskCreationTime) Len() int      { return len(b) }
func (b ByTaskCreationTime) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByTaskCreationTime) Less(i, j int) bool {
	return b[i].CreationTime < b[j].CreationTime
}
