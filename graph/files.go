package graph

import (
	"fmt"

	"github.com/bmeg/arachne/aql"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/ohsu-comp-bio/tes"
)

// File describes a file. Files maybe be of type "file" or "directory".
type File struct {
	URL  string
	Type tes.FileType
}

// MarshalAQL marshals the vertex into an arachne AQL vertex.
func (f *File) MarshalAQL() (*aql.Vertex, error) {
	if f.URL == "" {
		return nil, fmt.Errorf("can't marshal File: empty url")
	}
	return &aql.Vertex{
		Gid:   f.URL,
		Label: "Mortar.File",
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
