package v4config

import (
	"reflect"

	C "github.com/EHfive/v2ray-jsonschema/common"
	"github.com/iancoleman/orderedmap"
	JS "github.com/invopop/jsonschema"

	"github.com/v2fly/v2ray-core/v5/infra/conf/cfgcommon"
	"github.com/v2fly/v2ray-core/v5/infra/conf/rule"
	"github.com/v2fly/v2ray-core/v5/infra/conf/synthetic/dns"
	"github.com/v2fly/v2ray-core/v5/infra/conf/synthetic/router"
	v4 "github.com/v2fly/v2ray-core/v5/infra/conf/v4"
)

type CustomFakeDNSConfig struct{}
type CustomHostAddress struct{}

type CustomInboundConfig struct{}
type CustomOutboundConfig struct{}
type CustomTCPHeaderConfig struct{}
type CustomKCPHeaderConfig struct{}
type CustomMultiObservatoryItem struct{}
type CustomStrategyConfig struct{}

type CustomRouterRule struct {
	rule.RouterRule
	Domain     *cfgcommon.StringList  `json:"domain"`
	Domains    *cfgcommon.StringList  `json:"domains"`
	IP         *cfgcommon.StringList  `json:"ip"`
	Port       *cfgcommon.PortList    `json:"port"`
	Network    *cfgcommon.NetworkList `json:"network"`
	SourceIP   *cfgcommon.StringList  `json:"source"`
	SourcePort *cfgcommon.PortList    `json:"sourcePort"`
	User       *cfgcommon.StringList  `json:"user"`
	InboundTag *cfgcommon.StringList  `json:"inboundTag"`
	Protocols  *cfgcommon.StringList  `json:"protocol"`
	Attributes string                 `json:"attrs"`
}

func (CustomInboundConfig) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	return buildInOutBoundSchema(r, d, C.ToElemType((*v4.InboundDetourConfig)(nil)), inboundPairs)
}

func (CustomOutboundConfig) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	return buildInOutBoundSchema(r, d, C.ToElemType((*v4.OutboundDetourConfig)(nil)), outboundPairs)
}

func (CustomFakeDNSConfig) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	return C.BuildSingleOrArraySchema(r, d, C.ToElemType((*v4.FakeDNSPoolElementConfig)(nil)))
}

func (CustomHostAddress) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	return C.BuildSingleOrArraySchema(r, d, C.ToElemType((*C.CustomString)(nil)))
}

func (CustomTCPHeaderConfig) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	noneProps := orderedmap.New()
	noneProps.Set("type", &JS.Schema{Const: "none"})

	httpProps := orderedmap.New()
	httpProps.Set("type", &JS.Schema{Const: "http"})
	authS := r.RefOrReflectTypeToSchema(d, C.ToElemType((*v4.Authenticator)(nil)))

	return &JS.Schema{OneOf: []*JS.Schema{
		{Type: "object", Properties: noneProps},
		{AllOf: []*JS.Schema{
			{Type: "object", Properties: httpProps},
			authS,
		}},
	}}
}

func (CustomKCPHeaderConfig) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	types := []string{"none", "srtp", "utp", "wechat-video", "dtls", "wireguard"}
	s := &JS.Schema{}
	for _, name := range types {
		s.Enum = append(s.Enum, name)
	}
	props := orderedmap.New()
	props.Set("type", s)
	return &JS.Schema{Type: "object", Properties: props}
}

func (CustomMultiObservatoryItem) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	basicS := C.BuildBasicObjectSchema(r, d, C.ToElemType((*v4.MultiObservatoryItem)(nil)), []string{
		"type", "settings",
	})

	burstS := C.BuildOneOfItemSchema(r, d, "type", "settings", "burst", C.ToElemType((*v4.BurstObservatoryConfig)(nil)))
	defaultS := C.BuildOneOfItemSchema(r, d, "type", "settings", "default", C.ToElemType((*v4.ObservatoryConfig)(nil)))

	return &JS.Schema{
		AllOf: []*JS.Schema{
			basicS,
			{OneOf: []*JS.Schema{burstS, defaultS}},
		},
	}
}

func (CustomStrategyConfig) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	return C.BuildRouterStrategySchema(r, d, "type", "settings")
}

func alterField(t reflect.Type, f *reflect.StructField) bool {
	switch t {
	case C.ToElemType((*v4.TCPConfig)(nil)):
		if f.Name == "HeaderConfig" {
			f.Type = C.ToElemType((*CustomTCPHeaderConfig)(nil))
		}
	case C.ToElemType((*v4.KCPConfig)(nil)):
		if f.Name == "HeaderConfig" {
			f.Type = C.ToElemType((*CustomKCPHeaderConfig)(nil))
		}
	case C.ToElemType((*v4.QUICConfig)(nil)):
		if f.Name == "Header" {
			f.Type = C.ToElemType((*CustomKCPHeaderConfig)(nil))
		}
	}

	switch t {
	case C.ToElemType((*router.RouterConfig)(nil)):
		fallthrough
	case C.ToElemType((*router.RouterRulesConfig)(nil)):
		if f.Name == "RuleList" {
			f.Type = reflect.TypeOf(([]CustomRouterRule)(nil))
		}
	}

	matchType := f.Type
	if matchType.Kind() == reflect.Ptr {
		matchType = f.Type.Elem()
	}
	switch matchType {
	case reflect.TypeOf((*v4.InboundDetourConfig)(nil)):
		f.Type = reflect.TypeOf((*CustomInboundConfig)(nil))
	case reflect.TypeOf((*v4.OutboundDetourConfig)(nil)):
		f.Type = reflect.TypeOf((*CustomOutboundConfig)(nil))
	case reflect.TypeOf(([]v4.InboundDetourConfig)(nil)):
		f.Type = reflect.TypeOf(([]CustomInboundConfig)(nil))
	case reflect.TypeOf(([]v4.OutboundDetourConfig)(nil)):
		f.Type = reflect.TypeOf(([]CustomOutboundConfig)(nil))

	case reflect.TypeOf(([]v4.MultiObservatoryItem)(nil)):
		f.Type = reflect.TypeOf(([]CustomMultiObservatoryItem)(nil))
	case C.ToElemType((*v4.FakeDNSConfig)(nil)):
		f.Type = C.ToElemType((*CustomFakeDNSConfig)(nil))
	case C.ToElemType((*router.StrategyConfig)(nil)):
		f.Type = C.ToElemType((*CustomStrategyConfig)(nil))
	case reflect.TypeOf((map[string]*dns.HostAddress)(nil)):
		f.Type = reflect.TypeOf((map[string]CustomHostAddress)(nil))
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
	t := C.ToElemType((*v4.Config)(nil))
	return r.ReflectFromType(t)
}
