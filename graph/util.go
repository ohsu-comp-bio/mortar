package graph

import (
  "fmt"
  "github.com/bmeg/arachne/aql"
  "github.com/golang/protobuf/proto"
  structpb "github.com/golang/protobuf/ptypes/struct"
  "github.com/golang/protobuf/jsonpb"
)

var mar = jsonpb.Marshaler{}

type Vertex interface {
  MarshalAQLVertex() (*aql.Vertex, error)
}

type Edge interface {
  MarshalAQLEdge() (*aql.Edge, error)
}

type Client struct {
  *aql.Client
  Graph string
}
// TODO finish wrapping client

func (c *Client) AddVertex(v Vertex) error {
  av, err := v.MarshalAQLVertex()
  if err != nil {
    return err
  }
  return c.Client.AddVertex(c.Graph, *av)
}
func (c *Client) AddEdge(e Edge) error {
  ae, err := e.MarshalAQLEdge()
  if err != nil {
    return err
  }
  return c.Client.AddEdge(c.Graph, *ae)
}

func (c *Client) GetVertex(id string) (*aql.Vertex, error) {
  return c.Client.GetVertex(c.Graph, id)
}

func Marshal(msg proto.Message) (*structpb.Struct, error) {
  s, err := mar.MarshalToString(msg)
  if err != nil {
    return nil, err
  }

  st := &structpb.Struct{}
  err = jsonpb.UnmarshalString(s, st)
  if err != nil {
    return nil, err
  }

  return st, nil
}

func Unmarshal(st *structpb.Struct, msg proto.Message) error {
  b, err := mar.MarshalToString(st)
  if err != nil {
    return err
  }

  err = jsonpb.UnmarshalString(b, msg)
  if err != nil {
    return err
  }
  return nil
}

func NewEdge(label string, from, to Vertex) (*aql.Edge, error) {
  if label == "" {
    return nil, fmt.Errorf("can't create edge: empty label")
  }
  if from == nil {
    return nil, fmt.Errorf("can't create edge: empty From vertex")
  }
  if to == nil {
    return nil, fmt.Errorf("can't create edge: empty To vertex")
  }
  fv, err := from.MarshalAQLVertex()
  if err != nil {
    return nil, err
  }
  tv, err := to.MarshalAQLVertex()
  if err != nil {
    return nil, err
  }
  return &aql.Edge{
    Gid: fv.Gid + "->" + tv.Gid,
    Label: label,
    From: fv.Gid,
    To: tv.Gid,
  }, nil
}
