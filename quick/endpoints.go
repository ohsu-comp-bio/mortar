package main

import (
	"sort"

	"github.com/ohsu-comp-bio/mortar/graph"
	"github.com/ohsu-comp-bio/tes"
)


// TODO ideally, these models would directly reflect the models in the graph.
type runStat struct {
	Name      string
	Total     int
	Complete  int
	Steps     map[string]*graph.Step
	Tasks     map[string]*tes.Task
	StepTasks map[string][]string
}

type wfStat struct {
	Name string
	// TODO there will be more here later.
}

type respData struct {
	// Why separate IDs from data you ask?
	// Because this data is rendered into a sparse table.
	// The workflow and run data maps, might not contain entries for all table cells.
	WorkflowIDs []string
	RunIDs      []string

	Workflows map[string]*wfStat
	Runs      map[string]*runStat
}

func getData(cli graph.Client, graphID string) *respData {
	// TODO it is quite annoying to figure out how to best construct this data
	//      in the absence of DB sorting and grouping

	d := &respData{
		Workflows: map[string]*wfStat{},
		Runs:      map[string]*runStat{},
	}

	// TODO did I mention how complicated this is? I want grouping
	wfIDs := map[string]bool{}
	runIDs := map[string]bool{}

	res, err := cli.Query(graphID).
		V().
		HasLabel("ktl.Run").As("run").
		Out("ktl.RunForWorkflow").As("workflow").
		Select("run", "workflow").
		Execute()

	if err != nil {
		log.Error("ERR", err)
	}

	for row := range res {
		runv := row.Row[0].GetVertex()
		wfv := row.Row[1].GetVertex()
		d.Workflows[wfv.Gid] = &wfStat{wfv.Gid}
		wfIDs[wfv.Gid] = true
		runIDs[runv.Gid] = true
	}

	for wfid, _ := range wfIDs {
		d.WorkflowIDs = append(d.WorkflowIDs, wfid)
	}

	for rid, _ := range runIDs {
		d.RunIDs = append(d.RunIDs, rid)
		d.Runs[rid] = getRun(cli, graphID, rid)
	}

	// TODO think of a more meaningful sort field.
	//      is sorting by ID the right order in all cases?
	sort.Strings(d.WorkflowIDs)
	sort.Strings(d.RunIDs)

	return d
}

func getRun(cli graph.Client, graphID, runID string) *runStat {
	d := &runStat{
		Name: runID,
		Steps:     map[string]*graph.Step{},
		Tasks:     map[string]*tes.Task{},
		StepTasks: map[string][]string{},
	}

	// Get all steps
	res, err := cli.Query(graphID).
		V(runID).
		Out("ktl.RunForWorkflow").
		In("ktl.StepInWorkflow").
		Execute()

	if err != nil {
		log.Error("ERR", err)
	}
	for row := range res {
		v := row.Value.GetVertex()
		s := &graph.Step{Task: &tes.Task{}}
		err := graph.Unmarshal(v.Data, s)
		if err != nil {
			log.Error("Err", err)
			continue
		}
		d.Steps[v.Gid] = s
	}

	// Get steps with a task
	res, err = cli.Query(graphID).
		V(runID).
		Out("ktl.RunForWorkflow").
		In("ktl.StepInWorkflow").As("step").
		In("ktl.TaskForStep").As("task").
		Out("ktl.TaskForRun").
		HasId(runID).
		Select("step", "task").
		Execute()

	if err != nil {
		log.Error("ERR", err)
	}
	for row := range res {
		step := row.Row[0].GetVertex()
		task := &tes.Task{}
		err := graph.Unmarshal(row.Row[1].GetVertex().Data, task)
		if err != nil {
			panic(err)
		}
		d.StepTasks[step.Gid] = append(d.StepTasks[step.Gid], task.Id)
		d.Tasks[task.Id] = task
	}

	for sid, _ := range d.Steps {
		tids := d.StepTasks[sid]

		if len(tids) == 0 {
			continue
		}

		tasks := []*tes.Task{}
		for _, tid := range tids {
			tasks = append(tasks, d.Tasks[tid])
		}
		sort.Sort(ByCreationTime(tasks))
		latest := tasks[len(tasks)-1]

		if latest.State == tes.Complete {
			d.Complete++
		}
	}

	d.Total = len(d.Steps)
	log.Info("status", "run", runID, "total", d.Total, "complete", d.Complete)
	return d
}
