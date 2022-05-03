package common

import (
	"github.com/iancoleman/orderedmap"
	JS "github.com/invopop/jsonschema"
)

type CustomAny struct{}
type CustomString struct{}
type CustomStringList struct{}
type CustomPortRange struct{}

func (CustomString) JSONSchema() *JS.Schema {
	return &JS.Schema{Type: "string"}
}

func (CustomStringList) JSONSchema() *JS.Schema {
	return &JS.Schema{OneOf: []*JS.Schema{
		{Type: "string"},
		{Type: "array", Items: &JS.Schema{Type: "string"}},
	}}
}

func (CustomPortRange) JSONSchema() *JS.Schema {
	props := orderedmap.New()
	props.Set("from", &JS.Schema{Type: "integer"})
	props.Set("to", &JS.Schema{Type: "integer"})
	return &JS.Schema{OneOf: []*JS.Schema{
		{Type: "integer"},
		{Type: "object", Properties: props},
	}}
}
