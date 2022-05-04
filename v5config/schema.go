package v5config

import (
	"log"
	"reflect"

	C "github.com/EHfive/v2ray-jsonschema/common"
	"github.com/iancoleman/orderedmap"
	JS "github.com/invopop/jsonschema"

	"github.com/v2fly/v2ray-core/v5/app/router"
	"github.com/v2fly/v2ray-core/v5/infra/conf/v5cfg"

	_ "github.com/v2fly/v2ray-core/v5/main/distro/all"
)

type CustomInboundConfig struct{}
type CustomOutboundConfig struct{}
type CustomStreamSettings struct{}
type CustomServices struct{}
type CustomBalancingRule struct{}

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
		"grpc", "kcp", "tcp", "ws",
	})
	securitySList := C.BuildConditionalSchemaList(r, d, "security", "securitySettings", "security", []string{
		"tls",
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

func alterField(t reflect.Type, f *reflect.StructField) bool {
	switch t {
	case C.ToElemType((*v5cfg.RootConfig)(nil)):
		switch f.Name {
		case "LogConfig":
			f.Type = C.LoadTypeByAlias("service", "log")
		case "DNSConfig":
			f.Type = C.LoadTypeByAlias("service", "dns")
		case "RouterConfig":
			f.Type = C.LoadTypeByAlias("service", "router")
		case "Inbounds":
			f.Type = reflect.TypeOf(([]CustomInboundConfig)(nil))
		case "Outbounds":
			f.Type = reflect.TypeOf(([]CustomOutboundConfig)(nil))
		case "Services":
			f.Type = C.ToElemType((*CustomServices)(nil))
		}
	}

	matchType := f.Type
	if matchType.Kind() == reflect.Ptr {
		matchType = f.Type.Elem()
	}
	switch matchType {
	case C.ToElemType((*v5cfg.StreamConfig)(nil)):
		f.Type = C.ToElemType((*CustomStreamSettings)(nil))
	case reflect.TypeOf(([]*router.BalancingRule)(nil)):
		f.Type = reflect.TypeOf(([]*CustomBalancingRule)(nil))
	}

	if t.Name() == "router.RoutingRule" {
		log.Println(t, f)
	}

	return C.DefaultAlterField(t, f)
}

func customFields(t reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		skip := alterField(t, &f)
		if skip {
			continue
		}
		fields = append(fields, f)
	}
	return fields
}

func JSONSchema(r JS.Reflector) *JS.Schema {
	r.CustomFields = customFields
	t := C.ToElemType((*v5cfg.RootConfig)(nil))
	return r.ReflectFromType(t)
}
