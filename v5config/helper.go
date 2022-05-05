package v5config

import (
	"reflect"

	C "github.com/EHfive/v2ray-jsonschema/common"
	JS "github.com/invopop/jsonschema"
)

func buildInOutBoundSchema(r *JS.Reflector, d JS.Definitions, t reflect.Type, interfaceType string, protocols []string) *JS.Schema {
	idKey, configKey := "protocol", "settings"
	allOf := C.BuildConditionalSchemaList(r, d, idKey, configKey, interfaceType, protocols)
	allOf = append(allOf, C.BuildBasicObjectSchema(r, d, t, []string{idKey, configKey}))
	return &JS.Schema{AllOf: allOf, Required: []string{idKey}}
}
