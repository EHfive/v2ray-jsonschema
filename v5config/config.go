package v5config

import (
	"encoding/json"
	"reflect"

	"github.com/invopop/jsonschema"

	"github.com/v2fly/v2ray-core/v5/app/browserforwarder"
	"github.com/v2fly/v2ray-core/v5/app/dns"
	"github.com/v2fly/v2ray-core/v5/app/log"
	"github.com/v2fly/v2ray-core/v5/app/observatory"
	"github.com/v2fly/v2ray-core/v5/app/observatory/burst"
	"github.com/v2fly/v2ray-core/v5/app/policy"
	"github.com/v2fly/v2ray-core/v5/app/router"
	"github.com/v2fly/v2ray-core/v5/app/stats"
	"github.com/v2fly/v2ray-core/v5/infra/conf/v5cfg"
)

type V5Config struct {
	LogConfig    log.Config              `json:"log"`
	DNSConfig    dns.Config              `json:"dns"`
	RouterConfig router.SimplifiedConfig `json:"router"`
	Inbounds     []v5cfg.InboundConfig   `json:"inbounds"`
	Outbounds    []v5cfg.OutboundConfig  `json:"outbounds"`
	Services     V5Services              `json:"services"`
	Extensions   []json.RawMessage       `json:"extension"`
}

type V5Services struct {
	Browser          browserforwarder.Config `json:"browser"`
	Policy           policy.Config           `json:"policy"`
	Stats            stats.Config            `json:"stats"`
	BgObservatory    observatory.Config      `json:"backgroundObservatory"`
	BurstObservatory burst.Config            `json:"burstObservatory"`
}

func AdditionalFields(t reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fields = append(fields, f)
	}
	return fields
}

func JSONSchema(r jsonschema.Reflector) *jsonschema.Schema {
	r.AdditionalFields = AdditionalFields
	t := reflect.TypeOf(&V5Config{})
	return r.ReflectFromType(t)
}
