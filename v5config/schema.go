package v5config

import (
	"reflect"

	C "github.com/EHfive/v2ray-jsonschema/common"
	"github.com/iancoleman/orderedmap"
	JS "github.com/invopop/jsonschema"

	"github.com/v2fly/v2ray-core/v5/app/dns"
	"github.com/v2fly/v2ray-core/v5/app/router"
	"github.com/v2fly/v2ray-core/v5/infra/conf/v5cfg"
	"github.com/v2fly/v2ray-core/v5/transport/internet/tls/utls"

	_ "github.com/v2fly/v2ray-core/v5/main/distro/all"
)

type CustomInboundConfig struct{}
type CustomOutboundConfig struct{}
type CustomStreamSettings struct{}
type CustomServices struct{}
type CustomBalancingRule struct{}
type CustomUTLSImitate struct{}

func (CustomInboundConfig) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	return buildInOutBoundSchema(r, d, C.ToElemType((*v5cfg.InboundConfig)(nil)), "inbound", []string{
		"dokodemo-door",
		"http",
		"shadowsocks",
		"socks",
		"trojan",
		"vless",
		"vliteu",
		"vmess",
	})
}

func (CustomOutboundConfig) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	return buildInOutBoundSchema(r, d, C.ToElemType((*v5cfg.OutboundConfig)(nil)), "outbound", []string{
		"blackhole",
		"dns",
		"freedom",
		"http",
		"loopback",
		"shadowsocks",
		"socks",
		"trojan",
		"vless",
		"vliteu",
		"vmess",
	})
}

func (CustomStreamSettings) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	basicS := C.BuildBasicObjectSchema(r, d, C.ToElemType((*v5cfg.StreamConfig)(nil)), []string{
		"transport", "transportSettings", "security", "securitySettings",
	})
	transportSList := C.BuildConditionalSchemaList(r, d, "transport", "transportSettings", "transport", []string{
		"grpc", "kcp", "tcp", "quic", "ws",
	})
	securitySList := C.BuildConditionalSchemaList(r, d, "security", "securitySettings", "security", []string{
		"tls", "utls",
	})

	var allOf []*JS.Schema
	allOf = append(allOf, basicS)
	allOf = append(allOf, transportSList...)
	allOf = append(allOf, securitySList...)
	return &JS.Schema{AllOf: allOf}
}

func (CustomServices) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	services := []string{
		"backgroundObservatory",
		"browser",
		"burstObservatory",
		"commander",
		"fakeDns",
		"fakeDnsMulti",
		"instman",
		"policy",
		"restfulapi",
		"reverse",
		"stats",
	}

	props := orderedmap.New()
	for _, name := range services {
		props.Set(name, r.RefOrReflectTypeToSchema(d, C.LoadTypeByAlias("service", name)))
	}

	return &JS.Schema{
		Type:       "object",
		Properties: props,
	}
}

func (CustomBalancingRule) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	allOf := C.BuildRouterStrategySchemaList(r, d, "strategy", "strategySettings")
	basicS := C.BuildBasicObjectSchema(r, d, C.ToElemType((*router.BalancingRule)(nil)), []string{
		"strategy", "strategySettings",
	})
	allOf = append(allOf, basicS)
	return &JS.Schema{AllOf: allOf}
}

var imitateList = []string{
	"randomized",
	"randomizedalpn",
	"randomizednoalpn",
	"firefox_auto",
	"firefox_55",
	"firefox_56",
	"firefox_63",
	"firefox_65",
	"firefox_99",
	"firefox_102",
	"firefox_105",
	"chrome_auto",
	"chrome_58",
	"chrome_62",
	"chrome_70",
	"chrome_72",
	"chrome_83",
	"chrome_87",
	"chrome_96",
	"chrome_100",
	"chrome_102",
	"ios_auto",
	"ios_11_1",
	"ios_12_1",
	"ios_13",
	"ios_14",
	"android_11_okhttp",
	"edge_auto",
	"edge_85",
	"edge_106",
	"safari_auto",
	"safari_16_0",
	"360_auto",
	"360_7_5",
	"360_11_0",
	"qq_auto",
	"qq_11_1",
}

func (CustomUTLSImitate) JSONSchema() *JS.Schema {
	return C.BuildEnumSchema(imitateList)
}

var replaceFieldTypePairs []C.ReplaceFieldTypePair = []C.ReplaceFieldTypePair{
	{(*utls.Config)(nil), "Imitate", (*CustomUTLSImitate)(nil)},
	{(*dns.SimplifiedConfig)(nil), "DomainMatcher", (*C.CustomDNSDomainMatcher)(nil)},
	{(*router.SimplifiedRoutingRule)(nil), "DomainMatcher", (*C.CustomDNSDomainMatcher)(nil)},
}

var replaceTypePairs []C.ReplaceTypePair = []C.ReplaceTypePair{
	{(*v5cfg.InboundConfig)(nil), (*CustomInboundConfig)(nil)},
	{(*v5cfg.OutboundConfig)(nil), (*CustomOutboundConfig)(nil)},
	{(*v5cfg.StreamConfig)(nil), (*CustomStreamSettings)(nil)},
	{(*router.BalancingRule)(nil), (*CustomBalancingRule)(nil)},
}

func alterField(t reflect.Type, f *reflect.StructField) bool {
	switch t {
	case C.ToElemType((*v5cfg.RootConfig)(nil)):
		switch f.Name {
		case "LogConfig":
			f.Type = C.LoadTypeByAlias("service", "log")
			return false
		case "DNSConfig":
			f.Type = C.LoadTypeByAlias("service", "dns")
			return false
		case "RouterConfig":
			f.Type = C.LoadTypeByAlias("service", "router")
			return false
		case "Services":
			f.Type = C.ToElemType((*CustomServices)(nil))
			return false
		}
	}

	if newF, ok := C.ReplaceFieldTypeElemByPairs(t, *f, replaceFieldTypePairs); ok {
		f.Type = newF
		return false
	}

	if newF, ok := C.ReplaceTypeElemByPairs(f.Type, replaceTypePairs); ok {
		f.Type = newF
		return false
	}

	return C.DefaultAlterField(t, f)
}

func customFields(t reflect.Type) (fields []reflect.StructField) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		skip := alterField(t, &f)
		if skip {
			continue
		}
		fields = append(fields, f)
	}
	return
}

func JSONSchema(r JS.Reflector) *JS.Schema {
	r.CustomFields = customFields
	t := C.ToElemType((*v5cfg.RootConfig)(nil))
	s := r.ReflectFromType(t)
	return C.DefaultPostFixSchema(s, "jsonv5")
}
