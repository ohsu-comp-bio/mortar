package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/bmeg/arachne/aql"
	"github.com/golang/protobuf/proto"
	"github.com/ohsu-comp-bio/mortar/events"
	"github.com/ohsu-comp-bio/mortar/graph"
	"github.com/ohsu-comp-bio/tes"
	"github.com/rs/xid"
)

func main() {
	kw, err := events.NewKafkaWriter(events.KafkaConfig{
		Servers: []string{"localhost:9092"},
		Topic:   "funnel",
	})
	if err != nil {
		panic(err)
	}

	acli, err := aql.Connect("localhost:8202", true)
	if err != nil {
		panic(err)
	}

	cli := graph.Client{Client: &acli, Graph: "mortar"}
	b := &graph.Batch{}
	wg := sync.WaitGroup{}

	numWorkflows := 20
	numSamples := 50
	numSteps := 5

	for i := 0; i < numWorkflows; i++ {
		wf := &graph.Workflow{
			ID: fmt.Sprintf("wf-%.5d", i),
		}
		b.AddVertex(wf)

		for j := 0; j < numSteps; j++ {
			step := &graph.Step{
				ID: fmt.Sprintf("step-%.5d-%.5d", i, j),
			}
			e := graph.StepInWorkflow(step, wf)
			b.AddVertex(step)
			b.AddEdge(e)
		}

		for j := 0; j < numSamples; j++ {
			run := &graph.Run{
				ID:     fmt.Sprintf("run-%.5d-%.5d", i, j),
				Sample: fmt.Sprintf("sample-%.5d", j),
			}
			e := graph.RunForWorkflow(run, wf)
			b.AddVertex(run)
			b.AddEdge(e)

			wg.Add(1)
			go func(runID string, wfI int) {
				for k := 0; k < numSteps; k++ {
					simulateTask(kw, &tes.Task{
						Name: fmt.Sprintf("%s-step-%.5d", runID, k),
						Tags: map[string]string{
							"ktl.RunID":  runID,
							"ktl.StepID": fmt.Sprintf("step-%.5d-%.5d", wfI, k),
						},
					})
				}
				wg.Done()
			}(run.ID, i)
		}
	}

	err = cli.AddBatch(b)
	if err != nil {
		panic(err)
	}
	wg.Wait()
}

func simulateTask(kw *events.KafkaWriter, base *tes.Task) {
	task := proto.Clone(base).(*tes.Task)
	task.Id = xid.New().String()
	task.State = tes.Queued

	kw.Write(&events.Event{
		Id:   task.Id,
		Data: &events.Event_Task{Task: task},
		Type: events.Type_TASK_CREATED,
	})

	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)

	kw.Write(&events.Event{
		Id:   task.Id,
		Data: &events.Event_State{State: tes.Running},
		Type: events.Type_TASK_STATE,
	})

	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)

	kw.Write(&events.Event{
		Id:   task.Id,
		Data: &events.Event_State{State: tes.Complete},
		Type: events.Type_TASK_STATE,
	})
}
