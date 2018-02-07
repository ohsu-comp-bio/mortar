package graph

import (
	"testing"

	"github.com/ohsu-comp-bio/tes"
)

func TestNewEdge(t *testing.T) {
	task := &tes.Task{
		Id: "task-1",
	}
	a := &TaskVertex{task}
	b := &TagVertex{"foo", "zab"}
	e, err := NewEdge("Foo.Bar", a, b)
	if err != nil {
		t.Error("unexpected error", err)
	}
	if e.Gid != "task-1->foo:zab" {
		t.Error("unexpected edge id", e.Gid)
	}
	if e.Label != "Foo.Bar" {
		t.Error("unexpected edge label", e.Label)
	}
	if e.From != "task-1" {
		t.Error("unexpected edge from gid", e.From)
	}
	if e.To != "foo:zab" {
		t.Error("unexpected edge to gid", e.To)
	}
	_, err = NewEdge("Foo", a, nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = NewEdge("Foo", nil, b)
	if err == nil {
		t.Error("expected error")
	}
}
