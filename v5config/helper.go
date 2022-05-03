package v5config

import (
	"reflect"

	C "github.com/EHfive/v2ray-jsonschema/common"
	JS "github.com/invopop/jsonschema"
)

func buildInOutBoundSchema(r *JS.Reflector, d JS.Definitions, t reflect.Type, interfaceType string, protocols []string) *JS.Schema {
	return &JS.Schema{
		AllOf: []*JS.Schema{
			C.BuildBasicObjectSchema(r, d, t, []string{"protocol", "settings"}),
			C.BuildOneOfConfigsSchema(r, d, "protocol", "settings", interfaceType, protocols),
		},
	}
}
