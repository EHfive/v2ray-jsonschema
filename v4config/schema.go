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
	"github.com/v2fly/v2ray-core/v5/proxy/vless"
)

type CustomFakeDNSConfig struct{}
type CustomFakeDNSConfigExtend struct{}
type CustomHostAddress struct{}
type CustomNameServerConfig struct{}

// TODO: enums for queryStrategy, cacheStrategy and fallbackStrategy in NameServerConfig and DNSConfig

type CustomInboundConfig struct{}
type CustomOutboundConfig struct{}
type CustomTCPHeaderConfig struct{}
type CustomKCPHeaderConfig struct{}
type CustomMultiObservatoryItem struct{}
type CustomStrategyConfig struct{}
type CustomBlackholeConfigResponse struct{}
type CustomHTTPRemoteConfigUser struct{}
type CustomSocksRemoteConfigUser struct{}
type CustomVLessInOutboundConfigUser struct{}
type CustomVMessInOutboundConfigUser struct{}
type CustomRouterDomainStrategy struct{}
type CustomFreedomDomainStrategy struct{}
type CustomDNSQueryStrategy struct{}
type CustomDNSCacheStrategy struct{}
type CustomDNSFallbackStrategy struct{}
type CustomFallbackDest struct{}

type CustomRouterRule struct {
	rule.RouterRule
	Domain     *cfgcommon.StringList       `json:"domain"`
	Domains    *cfgcommon.StringList       `json:"domains"`
	IP         *cfgcommon.StringList       `json:"ip"`
	Port       *cfgcommon.PortList         `json:"port"`
	Network    *cfgcommon.NetworkList      `json:"network"`
	SourceIP   *cfgcommon.StringList       `json:"source"`
	SourcePort *cfgcommon.PortList         `json:"sourcePort"`
	User       *cfgcommon.StringList       `json:"user"`
	InboundTag *cfgcommon.StringList       `json:"inboundTag"`
	Protocols  *C.CustomRouterProtocolList `json:"protocol"`
	Attributes string                      `json:"attrs"`
}

func (CustomInboundConfig) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	return buildInOutBoundSchema(r, d, C.ToElemType((*v4.InboundDetourConfig)(nil)), inboundPairs)
}

func (CustomOutboundConfig) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	return buildInOutBoundSchema(r, d, C.ToElemType((*v4.OutboundDetourConfig)(nil)), outboundPairs)
}

func (CustomFakeDNSConfig) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	return C.BuildSingleOrArraySchema(r, d, C.ToElemType((*dns.FakeDNSPoolElementConfig)(nil)))
}

func (CustomFakeDNSConfigExtend) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	return &JS.Schema{OneOf: []*JS.Schema{
		{Type: "boolean"},
		C.SchemaFromPtr(r, d, (*CustomFakeDNSConfig)(nil)),
	}}
}

func (CustomHostAddress) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	return C.BuildSingleOrArraySchema(r, d, C.ToElemType((*C.CustomString)(nil)))
}

func (CustomNameServerConfig) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	return &JS.Schema{OneOf: []*JS.Schema{
		{Type: "string"},
		C.SchemaFromPtr(r, d, (*dns.NameServerConfig)(nil)),
	}}
}

func (CustomTCPHeaderConfig) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	noneProps := orderedmap.New()
	noneProps.Set("type", &JS.Schema{Const: "none"})
	noneS := &JS.Schema{
		If: &JS.Schema{Type: "object", Properties: noneProps},
	}

	httpProps := orderedmap.New()
	httpProps.Set("type", &JS.Schema{Const: "http"})
	authS := C.SchemaFromPtr(r, d, (*v4.Authenticator)(nil))
	httpS := &JS.Schema{
		If:   &JS.Schema{Type: "object", Properties: httpProps},
		Then: authS,
	}

	return &JS.Schema{AllOf: []*JS.Schema{noneS, httpS}}
}

func (CustomKCPHeaderConfig) JSONSchema() *JS.Schema {
	types := []string{"none", "srtp", "utp", "wechat-video", "dtls", "wireguard"}
	return buildObjectEnumSchema("type", types)
}

func (CustomMultiObservatoryItem) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	basicS := C.BuildBasicObjectSchema(r, d, C.ToElemType((*v4.MultiObservatoryItem)(nil)), []string{
		"type", "settings",
	})

	burstS := C.BuildConditionalItemSchema(r, d, "type", "settings", "burst", C.ToElemType((*v4.BurstObservatoryConfig)(nil)))
	defaultS := C.BuildConditionalItemSchema(r, d, "type", "settings", "default", C.ToElemType((*v4.ObservatoryConfig)(nil)))

	return &JS.Schema{
		AllOf: []*JS.Schema{
			basicS, burstS, defaultS,
		},
	}
}

func (CustomStrategyConfig) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	return &JS.Schema{AllOf: C.BuildRouterStrategySchemaList(r, d, "type", "settings")}
}

func (CustomBlackholeConfigResponse) JSONSchema() *JS.Schema {
	types := []string{"none", "http"}
	return buildObjectEnumSchema("type", types)
}

func (CustomHTTPRemoteConfigUser) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	return buildUserWithAccountSchema(r, d, (*v4.HTTPAccount)(nil))
}

func (CustomSocksRemoteConfigUser) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	return buildUserWithAccountSchema(r, d, (*v4.SocksAccount)(nil))
}

func (CustomVLessInOutboundConfigUser) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	return buildUserWithAccountSchema(r, d, (*vless.Account)(nil))
}

func (CustomVMessInOutboundConfigUser) JSONSchema2(r *JS.Reflector, d JS.Definitions) *JS.Schema {
	return buildUserWithAccountSchema(r, d, (*v4.VMessAccount)(nil))
}

func (CustomRouterDomainStrategy) JSONSchema() *JS.Schema {
	return C.BuildEnumSchema([]string{
		"AsIs",
		"AlwaysIP",
		"IPIfNonMatch",
		"IPOnDemand",
	})
}

func (CustomFreedomDomainStrategy) JSONSchema() *JS.Schema {
	return C.BuildEnumSchema([]string{
		"AsIs",
		"UseIP",
		"UseIPv4",
		"UseIPv6",
	})
}

func (CustomDNSQueryStrategy) JSONSchema() *JS.Schema {
	return C.BuildEnumSchema([]string{
		"UseIP",
		"UseIPv4",
		"UseIPv6",
	})
}

func (CustomDNSCacheStrategy) JSONSchema() *JS.Schema {
	return C.BuildEnumSchema([]string{
		"Enabled",
		"Disabled",
	})
}

func (CustomDNSFallbackStrategy) JSONSchema() *JS.Schema {
	return C.BuildEnumSchema([]string{
		"Enabled",
		"Disabled",
		"DisabledIfAnyMatch",
	})
}

func (CustomFallbackDest) JSONSchema() *JS.Schema {
	// port number, host:port or unix path, or "serve-ws-none"
	return &JS.Schema{OneOf: []*JS.Schema{
		{Type: "integer"},
		{Type: "string"},
	}}
}

var replaceFieldTypePairs []C.ReplaceFieldTypePair = []C.ReplaceFieldTypePair{
	{(*v4.TCPConfig)(nil), "HeaderConfig", (*CustomTCPHeaderConfig)(nil)},
	{(*v4.KCPConfig)(nil), "HeaderConfig", (*CustomKCPHeaderConfig)(nil)},
	{(*v4.QUICConfig)(nil), "Header", (*CustomKCPHeaderConfig)(nil)},
	{(*v4.BlackholeConfig)(nil), "Response", (*CustomBlackholeConfigResponse)(nil)},
	{(*v4.HTTPRemoteConfig)(nil), "Users", (*CustomHTTPRemoteConfigUser)(nil)},
	{(*v4.SocksRemoteConfig)(nil), "Users", (*CustomSocksRemoteConfigUser)(nil)},
	{(*v4.VLessInboundFallback)(nil), "Dest", (*CustomFallbackDest)(nil)},
	{(*v4.TrojanInboundFallback)(nil), "Dest", (*CustomFallbackDest)(nil)},
	{(*v4.VLessInboundConfig)(nil), "Clients", (*CustomVLessInOutboundConfigUser)(nil)},
	{(*v4.VLessOutboundVnext)(nil), "Users", (*CustomVLessInOutboundConfigUser)(nil)},
	{(*v4.VMessInboundConfig)(nil), "Users", (*CustomVMessInOutboundConfigUser)(nil)},
	{(*v4.VMessOutboundTarget)(nil), "Users", (*CustomVMessInOutboundConfigUser)(nil)},
	{(*router.RouterConfig)(nil), "RuleList", (*CustomRouterRule)(nil)},
	{(*router.RouterRulesConfig)(nil), "RuleList", (*CustomRouterRule)(nil)},
	{(*router.RouterConfig)(nil), "DomainStrategy", (*CustomRouterDomainStrategy)(nil)},
	{(*router.RouterRulesConfig)(nil), "DomainStrategy", (*CustomRouterDomainStrategy)(nil)},
	{(*v4.FreedomConfig)(nil), "DomainStrategy", (*CustomFreedomDomainStrategy)(nil)},
	{(*v4.OutboundDetourConfig)(nil), "DomainStrategy", (*CustomFreedomDomainStrategy)(nil)},
	{(*dns.DNSConfig)(nil), "QueryStrategy", (*CustomDNSQueryStrategy)(nil)},
	{(*dns.DNSConfig)(nil), "CacheStrategy", (*CustomDNSCacheStrategy)(nil)},
	{(*dns.DNSConfig)(nil), "FallbackStrategy", (*CustomDNSFallbackStrategy)(nil)},
	{(*dns.NameServerConfig)(nil), "QueryStrategy", (*CustomDNSQueryStrategy)(nil)},
	{(*dns.NameServerConfig)(nil), "CacheStrategy", (*CustomDNSCacheStrategy)(nil)},
	{(*dns.NameServerConfig)(nil), "FallbackStrategy", (*CustomDNSFallbackStrategy)(nil)},
	{(*dns.DNSConfig)(nil), "DomainMatcher", (*C.CustomDNSDomainMatcher)(nil)},
	{(*router.RouterConfig)(nil), "DomainMatcher", (*C.CustomDNSDomainMatcher)(nil)},
	{(*rule.RouterRule)(nil), "DomainMatcher", (*C.CustomDNSDomainMatcher)(nil)},
}

var replaceTypePairs []C.ReplaceTypePair = []C.ReplaceTypePair{
	{(*v4.InboundDetourConfig)(nil), (*CustomInboundConfig)(nil)},
	{(*v4.OutboundDetourConfig)(nil), (*CustomOutboundConfig)(nil)},
	{(*v4.MultiObservatoryItem)(nil), (*CustomMultiObservatoryItem)(nil)},
	{(*dns.FakeDNSConfig)(nil), (*CustomFakeDNSConfig)(nil)},
	{(*dns.FakeDNSConfigExtend)(nil), (*CustomFakeDNSConfigExtend)(nil)},
	{(*router.StrategyConfig)(nil), (*CustomStrategyConfig)(nil)},
	{(*dns.HostAddress)(nil), (*CustomHostAddress)(nil)},
	{(*dns.NameServerConfig)(nil), (*CustomNameServerConfig)(nil)},
}

func alterField(t reflect.Type, f *reflect.StructField) bool {
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
	s := r.ReflectFromType(t)
	return C.DefaultPostFixSchema(s, "jsonv4")
}
