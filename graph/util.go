package graph

import (
	"fmt"

	"github.com/bmeg/arachne/aql"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	structpb "github.com/golang/protobuf/ptypes/struct"
)

var mar = jsonpb.Marshaler{}

// Vertex describes a type of arachne AQL vertex.
type Vertex interface {
	MarshalAQL() (*aql.Vertex, error)
}

// Edge describes a type of arachne AQL edge.
type Edge interface {
	MarshalAQL() (*aql.Edge, error)
}

// Client wraps the arachne client with conveniences including:
// marshaling Vertex/Edge types using the MarshalAQL method,
// and being tied to a single graph.
type Client struct {
	*aql.Client
	Graph string
}

type Batch struct {
  Edges []Edge
  Verts []Vertex
}
func (b *Batch) AddVertex(v Vertex) {
  b.Verts = append(b.Verts, v)
}
func (b *Batch) AddEdge(e Edge) {
  b.Edges = append(b.Edges, e)
}

func (c *Client) AddBatch(b *Batch) error {
  // TODO what we really want is transaction semantics with rollback
  for _, v := range b.Verts {
    err := c.AddVertex(v)
    if err != nil {
      return fmt.Errorf("while adding vertex from batch: %s", err)
    }
  }
  for _, e := range b.Edges {
    err := c.AddEdge(e)
    if err != nil {
      return fmt.Errorf("while adding edge from batch: %s", err)
    }
  }
  return nil
}

// TODO finish wrapping client, and try to move this into arachne

// AddVertex adds a vertex to the graph.
func (c *Client) AddVertex(v Vertex) error {
	av, err := v.MarshalAQL()
	if err != nil {
		return err
	}
	return c.Client.AddVertex(c.Graph, *av)
}

// AddEdge adds an edge to the graph.
func (c *Client) AddEdge(e Edge) error {
	ae, err := e.MarshalAQL()
	if err != nil {
		return err
	}
	return c.Client.AddEdge(c.Graph, *ae)
}

// GetVertex gets a vertex by id.
func (c *Client) GetVertex(id string) (*aql.Vertex, error) {
	return c.Client.GetVertex(c.Graph, id)
}

// Marshal marshals a proto.Message into a structpb.Struct.
// Useful for preparing arachne requests.
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

// Unmarshal unmarshals a structpb.Struct into a proto.Message.
// Useful for unmarshaling arachne responses.
func Unmarshal(st *structpb.Struct, msg proto.Message) error {
	b, err := mar.MarshalToString(st)
	if err != nil {
		return err
	}

	return jsonpb.UnmarshalString(b, msg)
}

type edge struct {
	label    string
	from, to Vertex
}

func (e *edge) MarshalAQL() (*aql.Edge, error) {
	if e.label == "" {
		return nil, fmt.Errorf("can't create edge: empty label")
	}
	if e.from == nil {
		return nil, fmt.Errorf("can't create edge: empty From vertex")
	}
	if e.to == nil {
		return nil, fmt.Errorf("can't create edge: empty To vertex")
	}
	fv, err := e.from.MarshalAQL()
	if err != nil {
		return nil, err
	}
	tv, err := e.to.MarshalAQL()
	if err != nil {
		return nil, err
	}
	return &aql.Edge{
		Gid:   fv.Gid + "->" + tv.Gid,
		Label: e.label,
		From:  fv.Gid,
		To:    tv.Gid,
	}, nil
}

// NewEdge creates an arachne edge with the given `label`
// between the given `from` and `to` vertices. The edge's GID is constructed
// using the `from` and `to` vertices GIDs.
func NewEdge(label string, from, to Vertex) Edge {
	return &edge{label, from, to}
}
