package graph

import (
  "github.com/bmeg/arachne/aql"
  structpb "github.com/golang/protobuf/ptypes/struct"
)

type VertexLabel string
const (
  Task VertexLabel = "Funnel.Task"
  Tag = "Funnel.Task.Tag"
  File = "Mortar.File"
  Image = "Funnel.Container.Image"
)

func NewTagVertex(key, value string) *aql.Vertex {
  return &aql.Vertex{
    Gid: key + ":" + value,
    Label: Tag,
    Data: &structpb.Struct{
      Fields: map[string]*structpb.Value{
        "key": {
          &structpb.Value_StringValue{
            StringValue: value,
          },
        },
      },
    },
  }
}

func NewFileVertex(url string) *aql.Vertex {
  return &aql.Vertex{
    Gid: url,
    Label: File,
    Data: &structpb.Struct{
      Fields: map[string]*structpb.Value{
        "url": {
          Kind: &structpb.Value_StringValue{
            StringValue: url,
          },
        },
      },
    },
  }
}

func NewImageVertex(name string) *aql.Vertex {
  return &aql.Vertex{
    Gid: name,
    Label: Image,
    Data: &structpb.Struct{
      Fields: map[string]*structpb.Value{
        "name": {
          Kind: &structpb.Value_StringValue{
            StringValue: name,
          },
        },
      },
    },
  }
}
