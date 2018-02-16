package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/bmeg/arachne/aql"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/mortar/graph"
	"github.com/ohsu-comp-bio/tes"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var log = logger.NewLogger("quick", logger.DefaultConfig())

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

func main() {

	server := "localhost:8202"
	graphID := "quick"

	acli, err := aql.Connect(server, true)
	if err != nil {
		panic(err)
	}

	err = acli.AddGraph(graphID)
	if err != nil {
		panic(err)
	}

	cli := graph.Client{Client: &acli, Graph: graphID}

	b := makeData()
	if err := cli.AddBatch(b); err != nil {
		fmt.Println("ERR", err)
	}

	// Prometheus metrics
	http.HandleFunc("/metrics", func(resp http.ResponseWriter, req *http.Request) {
		updateMetrics(cli, graphID)
		promhttp.Handler().ServeHTTP(resp, req)
	})

	// JSON data for run/workflow/step/etc status
	http.HandleFunc("/data.json", func(resp http.ResponseWriter, req *http.Request) {
		d := getData(cli, graphID)
		enc := json.NewEncoder(resp)
		enc.SetIndent("", "  ")
		enc.Encode(d)
	})

	// Root web application
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("build/web"))))

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		http.ServeFile(resp, req, "build/web/index.html")
	})

	log.Info("listening", "http://localhost:9653")
	http.ListenAndServe(":9653", nil)
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
		// TODO don't like having everything be essentially a generic vertex
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
		log.Info("step data", v.Data)
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

type ByCreationTime []*tes.Task

func (b ByCreationTime) Len() int {
	return len(b)
}
func (b ByCreationTime) Less(i, j int) bool {
	return b[i].CreationTime < b[j].CreationTime
}
func (b ByCreationTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func fmtRow(row []*aql.QueryResult) []string {
	o := []string{}
	for _, item := range row {
		switch el := item.Result.(type) {
		// TODO this type switch is not intuitive. Should be aql.Vertex/Edge
		case *aql.QueryResult_Vertex:
			o = append(o, fmt.Sprintf("V(%s, %s)", el.Vertex.Label, el.Vertex.Gid))
		case *aql.QueryResult_Edge:
			o = append(o, fmt.Sprintf("E(%s, %s)", el.Edge.Label, el.Edge.Gid))
		}
	}
	return o
}
