package graph

import (
	"fmt"

	"github.com/bmeg/arachne/aql"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/ohsu-comp-bio/tes"
)

// Task is a vertex describing a task.
type Task struct {
	*tes.Task
}

// MarshalAQL marshals the vertex into an arachne AQL vertex.
func (t *Task) MarshalAQL() (*aql.Vertex, error) {
	if t.Id == "" {
		return nil, fmt.Errorf("can't marshal tes.Task: empty ID")
	}
	d, err := Marshal(t.Task)
	if err != nil {
		return nil, fmt.Errorf("can't marshal tes.Task: %s", err)
	}
	return &aql.Vertex{
		Gid:   t.Id,
		Label: "tes.Task",
		Data:  d,
	}, nil
}

// Image describes a container image (e.g. docker),
// such as the image described by a tes Task.
type Image struct {
	Name string
}

// MarshalAQL marshals the vertex into an arachne AQL vertex.
func (i *Image) MarshalAQL() (*aql.Vertex, error) {
	if i.Name == "" {
		return nil, fmt.Errorf("can't marshal Image: empty name")
	}
	return &aql.Vertex{
		Gid:   i.Name,
		Label: "tes.Task.Image",
		Data: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"name": {
					Kind: &structpb.Value_StringValue{
						StringValue: i.Name,
					},
				},
			},
		},
	}, nil
}

func TaskRequestsImage(t *Task, i *Image) Edge {
	return NewEdge("tes.Task.RequestsImage", t, i)
}

func TaskRequestsInput(t *Task, f *File) Edge {
	return NewEdge("tes.Task.RequestsInput", t, f)
}

func TaskRequestsOutput(t *Task, f *File) Edge {
	return NewEdge("tes.Task.RequestsOutput", t, f)
}

func TaskUploadedOutput(t *Task, f *File) Edge {
	return NewEdge("tes.Task.UploadedOutput", t, f)
}
