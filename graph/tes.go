package graph

import (
	"fmt"

	"github.com/bmeg/arachne/aql"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/ohsu-comp-bio/tes"
)

// TaskVertex is a vertex describing a task.
type TaskVertex struct {
	*tes.Task
}

// MarshalAQLVertex marshals the vertex into an arachne AQL vertex.
func (t *TaskVertex) MarshalAQLVertex() (*aql.Vertex, error) {
	if t.Task.Id == "" {
		return nil, fmt.Errorf("can't marshal TaskVertex: empty ID")
	}
	d, err := Marshal(t.Task)
	if err != nil {
		return nil, fmt.Errorf("can't marshal TaskVertex: %s", err)
	}
	return &aql.Vertex{
		Gid:   t.Task.Id,
		Label: "TES.Task",
		Data:  d,
	}, nil
}

// TagVertex is a vertex describing a Tag,
// likely created by and linked to a TES task.
type TagVertex struct {
	Key, Value string
}

// MarshalAQLVertex marshals the vertex into an arachne AQL vertex.
func (t *TagVertex) MarshalAQLVertex() (*aql.Vertex, error) {
	if t.Key == "" {
		return nil, fmt.Errorf("can't marshal TagVertex: empty key")
	}
	return &aql.Vertex{
		Gid:   t.Key + ":" + t.Value,
		Label: "TES.Task.Tag",
		Data: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"key": {
					&structpb.Value_StringValue{
						StringValue: t.Key,
					},
				},
				"value": {
					&structpb.Value_StringValue{
						StringValue: t.Value,
					},
				},
			},
		},
	}, nil
}

// FileVertex describes a file. Files maybe be of type "file" or "directory".
type FileVertex struct {
	URL  string
	Type tes.FileType
}

// MarshalAQLVertex marshals the vertex into an arachne AQL vertex.
func (f *FileVertex) MarshalAQLVertex() (*aql.Vertex, error) {
	if f.URL == "" {
		return nil, fmt.Errorf("can't marshal FileVertex: empty url")
	}
	return &aql.Vertex{
		Gid:   f.URL,
		Label: "Mortar.File.URL",
		Data: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"url": {
					Kind: &structpb.Value_StringValue{
						StringValue: f.URL,
					},
				},
				"type": {
					Kind: &structpb.Value_StringValue{
						StringValue: f.Type.String(),
					},
				},
			},
		},
	}, nil
}

// ImageVertex describes a container image (e.g. docker),
// such as the image described by a TES Task.
type ImageVertex struct {
	Name string
}

// MarshalAQLVertex marshals the vertex into an arachne AQL vertex.
func (i *ImageVertex) MarshalAQLVertex() (*aql.Vertex, error) {
	if i.Name == "" {
		return nil, fmt.Errorf("can't marshal ImageVertex: empty name")
	}
	return &aql.Vertex{
		Gid:   i.Name,
		Label: "TES.Task.Image",
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

// HasTagEdge links a tag to a task.
type HasTagEdge struct {
	From *TaskVertex
	To   *TagVertex
}

// MarshalAQLEdge marshals the edge into an arachne AQL Edge.
func (e *HasTagEdge) MarshalAQLEdge() (*aql.Edge, error) {
	return NewEdge("TES.Task.HasTag", e.From, e.To)
}

// RequestsImageEdge links an image to a task.
type RequestsImageEdge struct {
	From *TaskVertex
	To   *ImageVertex
}

// MarshalAQLEdge marshals the edge into an arachne AQL Edge.
func (e *RequestsImageEdge) MarshalAQLEdge() (*aql.Edge, error) {
	return NewEdge("TES.Task.RequestsImage", e.From, e.To)
}

// RequestsInputEdge links a file to a task input description.
type RequestsInputEdge struct {
	From *TaskVertex
	To   *FileVertex
}

// MarshalAQLEdge marshals the edge into an arachne AQL Edge.
func (e *RequestsInputEdge) MarshalAQLEdge() (*aql.Edge, error) {
	return NewEdge("TES.Task.RequestsInput", e.From, e.To)
}

// RequestsOutputEdge links a file to a task output description.
type RequestsOutputEdge struct {
	From *TaskVertex
	To   *FileVertex
}

// MarshalAQLEdge marshals the edge into an arachne AQL Edge.
func (e *RequestsOutputEdge) MarshalAQLEdge() (*aql.Edge, error) {
	return NewEdge("TES.Task.RequestsOutput", e.From, e.To)
}

// UploadedOutputEdge links a task to an output file.
type UploadedOutputEdge struct {
	From *TaskVertex
	To   *FileVertex
}

// MarshalAQLEdge marshals the edge into an arachne AQL Edge.
func (e *UploadedOutputEdge) MarshalAQLEdge() (*aql.Edge, error) {
	return NewEdge("TES.Task.UploadedOutput", e.From, e.To)
}
