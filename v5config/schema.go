package v5config

import (
	"reflect"
	"strings"

	C "github.com/EHfive/v2ray-jsonschema/common"
	"github.com/iancoleman/orderedmap"
	JS "github.com/invopop/jsonschema"

	"github.com/v2fly/v2ray-core/v5/infra/conf/v5cfg"

	_ "github.com/v2fly/v2ray-core/v5/main/distro/all"
)

type CustomInboundConfig struct{}
type CustomOutboundConfig struct{}
type CustomStreamSettings struct{}
type CustomServices struct{}

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
	basic := buildBasicObjectSchema(r, d, C.ToElemType((*v5cfg.StreamConfig)(nil)), []string{
		"transport",
		"transportSettings",
		"security",
		"securitySettings",
	})

	transport := buildOneOfConfigsSchema(r, d, "transport", "transportSettings", "transport", []string{
		"grpc",
		"kcp",
		"tcp",
		"ws",
	})

	security := buildOneOfConfigsSchema(r, d, "security", "securitySettings", "security", []string{
		"tls",
	})

	return &JS.Schema{
		AllOf: []*JS.Schema{
			basic,
			transport,
			security,
		},
	}
}

func (CustomServices) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	services := []string{
		"browser",
		"policy",
		"stats",
		"backgroundObservatory",
		"burstObservatory",
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

func buildBasicObjectSchema(r *JS.Reflector, d JS.Definitions, t reflect.Type, excludes []string) *JS.Schema {
	s := r.RefOrReflectTypeToSchema(d, t)
	res := s
	if s.Ref != "" {
		// s = d[s.Ref]
		defName := strings.Replace(s.Ref, "#/$defs/", "", 1)
		s = d[defName]
	}
	for _, name := range excludes {
		s.Properties.Delete(name)
	}
	s.AdditionalProperties = JS.TrueSchema
	return res
}

func buildOneOfConfigsSchema(r *JS.Reflector, d JS.Definitions, idKey string, configKey string, interfaceType string, shortNames []string) *JS.Schema {
	var schemas []*JS.Schema
	for _, name := range shortNames {
		props := orderedmap.New()
		props.Set(idKey, &JS.Schema{Const: name})
		props.Set(configKey, r.RefOrReflectTypeToSchema(d, C.LoadTypeByAlias(interfaceType, name)))
		schema := &JS.Schema{
			AdditionalProperties: JS.TrueSchema,
			Properties:           props,
		}
		schemas = append(schemas, schema)
	}
	return &JS.Schema{
		OneOf: schemas,
	}
}

func buildInOutBoundSchema(r *JS.Reflector, d JS.Definitions, t reflect.Type, interfaceType string, protocols []string) *JS.Schema {
	return &JS.Schema{
		AllOf: []*JS.Schema{
			buildBasicObjectSchema(r, d, t, []string{"protocol", "settings"}),
			buildOneOfConfigsSchema(r, d, "protocol", "settings", interfaceType, protocols),
		},
	}
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
	case C.ToElemType((*v5cfg.InboundConfig)(nil)):
		switch f.Name {
		case "StreamSetting":
			f.Type = C.ToElemType((*CustomStreamSettings)(nil))
		}
	case C.ToElemType((*v5cfg.OutboundConfig)(nil)):
		switch f.Name {
		case "StreamSetting":
			f.Type = C.ToElemType((*CustomStreamSettings)(nil))
		}
	}

	return false
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
