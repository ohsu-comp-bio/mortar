package graph

import (
  "fmt"
  "github.com/bmeg/arachne/aql"
  "github.com/ohsu-comp-bio/tes"
  structpb "github.com/golang/protobuf/ptypes/struct"
)


type TaskVertex struct {
  *tes.Task
}

func (t *TaskVertex) MarshalAQLVertex() (*aql.Vertex, error) {
  if t.Task.Id == "" {
    return nil, fmt.Errorf("can't marshal TaskVertex: empty ID")
  }
  d, err := Marshal(t.Task)
  if err != nil {
    return nil, fmt.Errorf("can't marshal TaskVertex: %s", err)
  }
  return &aql.Vertex{
    Gid: t.Task.Id,
    Label: "Funnel.Task",
    Data: d,
  }, nil
}

type TagVertex struct {
  Key, Value string
}

func (t *TagVertex) MarshalAQLVertex() (*aql.Vertex, error) {
  if t.Key == "" {
    return nil, fmt.Errorf("can't marshal TagVertex: empty key")
  }
  return &aql.Vertex{
    Gid: t.Key + ":" + t.Value,
    Label: "Funnel.Task.Tag",
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

type FileVertex struct {
  Url string
  Type tes.FileType
}

func (f *FileVertex) MarshalAQLVertex() (*aql.Vertex, error) {
  if f.Url == "" {
    return nil, fmt.Errorf("can't marshal FileVertex: empty url")
  }
  return &aql.Vertex{
    Gid: f.Url,
    Label: "Mortar.File.Url",
    Data: &structpb.Struct{
      Fields: map[string]*structpb.Value{
        "url": {
          Kind: &structpb.Value_StringValue{
            StringValue: f.Url,
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

type ImageVertex struct {
  Name string
}

func (i *ImageVertex) MarshalAQLVertex() (*aql.Vertex, error) {
  if i.Name == "" {
    return nil, fmt.Errorf("can't marshal ImageVertex: empty name")
  }
  return &aql.Vertex{
    Gid: i.Name,
    Label: "Funnel.Task.Image",
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

type HasTagEdge struct {
  From *TaskVertex
  To *TagVertex
}
func (e *HasTagEdge) MarshalAQLEdge() (*aql.Edge, error) {
  return NewEdge("Funnel.Task.HasTag", e.From, e.To)
}

type RequestsImageEdge struct {
  From *TaskVertex
  To *ImageVertex
}

func (e *RequestsImageEdge) MarshalAQLEdge() (*aql.Edge, error) {
  return NewEdge("Funnel.Task.RequestsImage", e.From, e.To)
}

type RequestsInputEdge struct {
  From *TaskVertex
  To *FileVertex
}

func (e *RequestsInputEdge) MarshalAQLEdge() (*aql.Edge, error) {
  return NewEdge("Funnel.Task.RequestsInput", e.From, e.To)
}

type RequestsOutputEdge struct {
  From *TaskVertex
  To *FileVertex
}

func (e *RequestsOutputEdge) MarshalAQLEdge() (*aql.Edge, error) {
  return NewEdge("Funnel.Task.RequestsOutput", e.From, e.To)
}

type UploadedOutputEdge struct {
  From *TaskVertex
  To *FileVertex
}

func (e *UploadedOutputEdge) MarshalAQLEdge() (*aql.Edge, error) {
  return NewEdge("Funnel.Task.UploadedOutput", e.From, e.To)
}
