package common

import (
	"reflect"

	"github.com/iancoleman/orderedmap"
	JS "github.com/invopop/jsonschema"
	"github.com/v2fly/v2ray-core/v5/common/net"
)

type CustomPbAny struct{}
type CustomString struct{}
type CustomStringList struct{}
type CustomPortRange struct{}
type CustomNumber struct{}
type CustomNetworkList struct{}

func (CustomPbAny) JSONSchema() *JS.Schema {
	props := orderedmap.New()
	props.Set("typeUrl", &JS.Schema{Type: "string", Pattern: "^types\\.v2fly\\.org/"})
	props.Set("value", &JS.Schema{Type: "string", ContentEncoding: "base64"})
	return &JS.Schema{Type: "object", Properties: props}
}

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
	return &JS.Schema{OneOf: []*JS.Schema{
		{Type: "integer"},
		{Type: "string"},
	}}
}

func (CustomNumber) JSONSchema() *JS.Schema {
	return &JS.Schema{Type: "integer"}
}

func (CustomNetworkList) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	s := r.RefOrReflectTypeToSchema(d, reflect.TypeOf(([]net.Network)(nil)))
	return &JS.Schema{OneOf: []*JS.Schema{
		{Type: "string"},
		s,
	}}
}
