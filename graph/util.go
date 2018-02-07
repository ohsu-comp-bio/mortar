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
	MarshalAQLVertex() (*aql.Vertex, error)
}

// Edge describes a type of arachne AQL edge.
type Edge interface {
	MarshalAQLEdge() (*aql.Edge, error)
}

// Client wraps the arachne client with conveniences including:
// marshaling Vertex/Edge types using the MarshalAQLVertex/MarshalAQLEdge method,
// and being tied to a single graph.
type Client struct {
	*aql.Client
	Graph string
}

// TODO finish wrapping client, and try to move this into arachne

// AddVertex adds a vertex to the graph.
func (c *Client) AddVertex(v Vertex) error {
	av, err := v.MarshalAQLVertex()
	if err != nil {
		return err
	}
	return c.Client.AddVertex(c.Graph, *av)
}

// AddEdge adds an edge to the graph.
func (c *Client) AddEdge(e Edge) error {
	ae, err := e.MarshalAQLEdge()
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

// NewEdge creates an arachne edge with the given `label`
// between the given `from` and `to` vertices. The edge's GID is constructed
// using the `from` and `to` vertices GIDs.
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
		Gid:   fv.Gid + "->" + tv.Gid,
		Label: label,
		From:  fv.Gid,
		To:    tv.Gid,
	}, nil
}
