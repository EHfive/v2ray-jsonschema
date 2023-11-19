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
type CustomNetworkList struct{}
type CustomDNSDomainMatcher struct{}
type CustomRouterProtocol struct{}
type CustomRouterProtocolList struct{}
type CustomFreedomDomainStrategy struct{}

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

func (CustomNetworkList) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	s := r.RefOrReflectTypeToSchema(d, reflect.TypeOf(([]net.Network)(nil)))
	return &JS.Schema{OneOf: []*JS.Schema{
		{Type: "string"},
		s,
	}}
}

func (CustomDNSDomainMatcher) JSONSchema() *JS.Schema {
	return BuildEnumSchema([]string{"linear", "mph"})
}

func (CustomRouterProtocol) JSONSchema() *JS.Schema {
	return BuildEnumSchema([]string{
		// from <https://pkg.go.dev/github.com/v2fly/v2ray-core/v5/common/session#Content>
		"tls", "quic", "dns",
		// from <https://pkg.go.dev/github.com/v2fly/v2ray-core/v5/app/dispatcher#SniffResult>
		"fakedns", "fakedns+others", "bittorrent", "http1", // "quic", "tls",
	})
}

func (CustomRouterProtocolList) JSONSchema() *JS.Schema {
	s := CustomRouterProtocol{}.JSONSchema()
	return &JS.Schema{OneOf: []*JS.Schema{
		s,
		{Type: "array", Items: s},
	}}
}

func (CustomFreedomDomainStrategy) JSONSchema() *JS.Schema {
	return BuildEnumSchema([]string{
		"AsIs",
		"UseIP",
		"UseIPv4",
		"UseIPv6",
	})
}
