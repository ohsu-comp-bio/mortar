package main

import (
  "fmt"
	"github.com/bmeg/arachne/aql"
	"github.com/ohsu-comp-bio/funnel/logger"
	"github.com/ohsu-comp-bio/mortar/graph"
	"github.com/ohsu-comp-bio/tes"
)

var log = logger.NewLogger("quick", logger.DefaultConfig())

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
  b := &graph.Batch{}

  wf1 := &graph.Workflow{"WF1"}

  s1 := &graph.Step{"S1", &tes.Task{
    Name: "S1",
    Description: "Mortar step example",
    Executors: []*tes.Executor{
      {
        Image: "alpine",
        Command: []string{"echo", "hello", "world"},
      },
    },
  }}


  s2 := &graph.Step{"S2", &tes.Task{
    Name: "S2",
    Description: "Mortar step 2 example",
    Executors: []*tes.Executor{
      {
        Image: "alpine",
        Command: []string{"md5sum", "example.txt"},
      },
    },
  }}

  s3 := &graph.Step{"S3", &tes.Task{
    Name: "S3",
    Description: "Mortar step 3 example",
    Executors: []*tes.Executor{
      {
        Image: "alpine",
        Command: []string{"cat"},
      },
    },
  }}

  r1t1 := &graph.Task{&tes.Task{
    Id: "run1-123",
    Name: "Run1 S1",
    Description: "Mortar step example",
    Executors: []*tes.Executor{
      {
        Image: "alpine",
        Command: []string{"echo", "hello", "world"},
      },
    },
  }}

  r1t2 := &graph.Task{&tes.Task{
    Id: "run1-124",
    Name: "Run1 S2",
    Description: "Mortar step 2 example",
    Executors: []*tes.Executor{
      {
        Image: "alpine",
        Command: []string{"md5sum", "example.txt"},
      },
    },
  }}

  r2t1 := &graph.Task{&tes.Task{
    Id: "run2-125",
    Name: "Run2 S1",
    Description: "Mortar step example",
    Executors: []*tes.Executor{
      {
        Image: "alpine",
        Command: []string{"echo", "hello", "world"},
      },
    },
  }}

  r2t2 := &graph.Task{&tes.Task{
    Id: "run2-126",
    Name: "Run2 S2",
    Description: "Mortar step 2 example",
    Executors: []*tes.Executor{
      {
        Image: "alpine",
        Command: []string{"md5sum", "example.txt"},
      },
    },
  }}

  r2t3 := &graph.Task{&tes.Task{
    Id: "run2-127",
    Name: "Run2 S3",
    Description: "Mortar step 3 example",
    Executors: []*tes.Executor{
      {
        Image: "alpine",
        Command: []string{"md5sum", "example.txt"},
      },
    },
  }}

  r1 := &graph.Run{"Run 1"}
  r2 := &graph.Run{"Run 2"}

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

  if err := cli.AddBatch(b); err != nil {
    fmt.Println("ERR", err)
  }

  // Begin queries

  // Start simple, just retrieve a run.
  res, err := cli.Query(graphID).V("Run 1").Execute()
  if err != nil {
    fmt.Println("ERR", err)
  }
  for row := range res {
    fmt.Println("ROW", row)
    fmt.Println("VALUE", row.Value)
    r1 := row.Value.GetVertex()
    fmt.Println("RUN", r1)
  }

  // Get all the tasks connected to a run.
  // Note that Select() changes the return type to use the res.Row field.
  res, err = cli.Query(graphID).
    V("Run 1").As("run").
    In("ktl.TaskForRun").As("task").
    Select("run", "task").Execute()

  if err != nil {
    fmt.Println("ERR", err)
  }
  for row := range res {
    log.Info("q2", fmtRow(row.Row))
  }

  // For a given run, get all the tasks, then all the steps those tasks
  // are for.
  res, err = cli.Query(graphID).
    V("Run 1").As("run").
    In("ktl.TaskForRun").As("task").
    Out("ktl.TaskForStep").As("step").
    Select("run", "task", "step").Execute()

  if err != nil {
    fmt.Println("ERR", err)
  }
  for row := range res {
    log.Info("q3", fmtRow(row.Row))
  }

  // Get the workflow that a run is running.
  res, err = cli.Query(graphID).
    V("Run 1").
    Out("ktl.RunForWorkflow").
    Execute()

  if err != nil {
    fmt.Println("ERR", err)
  }

  for row := range res {
    log.Info("q4", row.Value)
  }


  // DEBUG get all ktl.StepInWorkfow edges
  res, err = cli.Query(graphID).
    E().
    HasLabel("ktl.StepInWorkflow").
    Execute()

  if err != nil {
    log.Error("ERR", err)
  }

  for row := range res {
    log.Info("q5", row)
  }

  // Get all the steps in the workflow that is running.
  res, err = cli.Query(graphID).
    V("Run 1").
    Out("ktl.RunForWorkflow").
    In("ktl.StepInWorkflow").
    // TODO spend 10-20 minutes debugging this line, which has a typo "Workfow"
    //      very subtle. want help from the client code and/or server schema.
    //In("ktl.StepInWorkfow").
    Execute()

  if err != nil {
    log.Error("ERR", err)
  }
  for row := range res {
    v := row.Value.GetVertex()
    log.Info("all steps in a running workflow", v.Label, v.Gid)
  }

  // Get all the steps in the workflow that is running, and get the
  // tasks for each step in the run.
  // Get all the steps in the workflow that is running.
  res, err = cli.Query(graphID).
    V("Run 1").
    Out("ktl.RunForWorkflow").
    In("ktl.StepInWorkflow").As("step").
    In("ktl.TaskForStep").As("task").
    // TODO how is this even working right?
    //      the select is in the wrong place.
    Select("step", "task").
    Out("ktl.TaskForRun").
    HasId("Run 1").
    Execute()

  if err != nil {
    log.Error("ERR", err)
  }
  for row := range res {
    log.Info("q7", fmtRow(row.Row))
  }


  // Now the same for run 2
  res, err = cli.Query(graphID).
    V("Run 2").
    Out("ktl.RunForWorkflow").
    In("ktl.StepInWorkflow").As("step").
    In("ktl.TaskForStep").As("task").
    Out("ktl.TaskForRun").
    HasId("Run 2").
    Select("step", "task").
    Execute()

  if err != nil {
    log.Error("ERR", err)
  }
  for row := range res {
    log.Info("steps with tasks for a run", fmtRow(row.Row))
  }

  // TODO the query above returns only the running steps.
  //      would be nice to include the tasks only if the exist, null otherwise



  // Now combine all steps + steps with tasks
  //
  // For a given run, how complete is the run?
  // Get all the steps in the running workflow, and determine whether they
  // are all complete.
  steps := map[string]*aql.Vertex{}
  res, err = cli.Query(graphID).
    V("Run 2").
    Out("ktl.RunForWorkflow").
    In("ktl.StepInWorkflow").
    Execute()

  if err != nil {
    log.Error("ERR", err)
  }
  for row := range res {
    v := row.Value.GetVertex()
    steps[v.Gid] = v
    log.Info("steps in workflow", v)
  }

  stepTasks := map[string][]*aql.Vertex{}
  res, err = cli.Query(graphID).
    V("Run 2").
    Out("ktl.RunForWorkflow").
    In("ktl.StepInWorkflow").As("step").
    In("ktl.TaskForStep").As("task").
    Out("ktl.TaskForRun").
    HasId("Run 2").
    Select("step", "task").
    Execute()

  if err != nil {
    log.Error("ERR", err)
  }
  for row := range res {
    step := row.Row[0].GetVertex()
    task := row.Row[1].GetVertex()
    stepTasks[step.Gid] = append(stepTasks[step.Gid], task)
  }

  running := 0
  for sid, _ := range steps {
    tasks := stepTasks[sid]
    if len(tasks) > 1 {
      panic("unhandled case, multiple tasks for a running step")
    }
    if len(tasks) == 1 {
      running++
    }
  }

  log.Info("status", "steps", len(steps), "running", running)


  // Try doing the above analysis for all runs at once
  res, err = cli.Query(graphID).
    V().
    HasLabel("ktl.Run").As("run").
    Out("ktl.RunForWorkflow").As("workflow").
    In("ktl.StepInWorkflow").As("step").
    Select("run", "workflow", "step").
    Execute()

  if err != nil {
    log.Error("ERR", err)
  }
  for row := range res {
    log.Info("q9", fmtRow(row.Row))
  }

  // TODO where() or filterValues() is  needed for part 2, because the
  //      tasks need to be filtered on the current run.



  res, err = cli.Query(graphID).
    V().
    HasLabel("ktl.Run").As("run").
    Execute()

  if err != nil {
    log.Error("ERR", err)
  }
  for row := range res {
    v := row.Value.GetVertex()
    getRun(cli, graphID, v.Gid)
  }
}

func getRun(cli graph.Client, graphID, runID string) {
  steps := map[string]*aql.Vertex{}
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
    steps[v.Gid] = v
  }

  stepTasks := map[string][]*aql.Vertex{}
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
    task := row.Row[1].GetVertex()
    stepTasks[step.Gid] = append(stepTasks[step.Gid], task)
  }

  running := 0
  for sid, _ := range steps {
    tasks := stepTasks[sid]
    if len(tasks) > 1 {
      panic("unhandled case, multiple tasks for a running step")
    }
    if len(tasks) == 1 {
      running++
    }
  }

  log.Info("status", "run", runID, "steps", len(steps), "running", running)
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
