package v4config

import (
	"reflect"

	C "github.com/EHfive/v2ray-jsonschema/common"
	JS "github.com/invopop/jsonschema"

	v4 "github.com/v2fly/v2ray-core/v5/infra/conf/v4"
)

type ConfigPair struct {
	Name      string
	Interface interface{}
}

var inboundPairs = []ConfigPair{
	{"dokodemo-door", ((*v4.DokodemoConfig)(nil))},
	{"http", ((*v4.HTTPServerConfig)(nil))},
	{"shadowsocks", ((*v4.ShadowsocksServerConfig)(nil))},
	{"socks", ((*v4.SocksServerConfig)(nil))},
	{"vless", ((*v4.VLessInboundConfig)(nil))},
	{"vmess", ((*v4.VMessInboundConfig)(nil))},
	{"trojan", ((*v4.TrojanServerConfig)(nil))},
}

var outboundPairs = []ConfigPair{
	{"blackhole", ((*v4.BlackholeConfig)(nil))},
	{"freedom", ((*v4.FreedomConfig)(nil))},
	{"http", ((*v4.HTTPClientConfig)(nil))},
	{"shadowsocks", ((*v4.ShadowsocksClientConfig)(nil))},
	{"socks", ((*v4.SocksClientConfig)(nil))},
	{"vless", ((*v4.VLessOutboundConfig)(nil))},
	{"vmess", ((*v4.VMessOutboundConfig)(nil))},
	{"trojan", ((*v4.TrojanClientConfig)(nil))},
	{"dns", ((*v4.DNSOutboundConfig)(nil))},
	{"loopback", ((*v4.LoopbackConfig)(nil))},
}

func buildInOutBoundSchema(r *JS.Reflector, d JS.Definitions, t reflect.Type, configPairs []ConfigPair) *JS.Schema {
	idKey := "protocol"
	configKey := "settings"

	var schemas []*JS.Schema
	for _, pair := range configPairs {
		s := C.BuildOneOfItemSchema(r, d, idKey, configKey, pair.Name, C.ToElemType(pair.Interface))
		schemas = append(schemas, s)
	}
	return &JS.Schema{
		AllOf: []*JS.Schema{
			C.BuildBasicObjectSchema(r, d, t, []string{idKey, configKey}),
			{OneOf: schemas},
		},
	}
}
