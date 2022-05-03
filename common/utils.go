package common

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/iancoleman/orderedmap"
	JS "github.com/invopop/jsonschema"
	"github.com/stoewer/go-strcase"

	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/infra/conf/cfgcommon"
	"github.com/v2fly/v2ray-core/v5/infra/conf/cfgcommon/duration"

	"google.golang.org/protobuf/types/known/anypb"

	"github.com/v2fly/v2ray-core/v5/common/environment/envctx"
	"github.com/v2fly/v2ray-core/v5/common/environment/envimpl"
	"github.com/v2fly/v2ray-core/v5/common/registry"
)

func NewDefaultReflector() JS.Reflector {
	return JS.Reflector{
		RequiredFromJSONSchemaTags: true,
		AllowAdditionalProperties:  true,
		Namer: func(t reflect.Type) string {
			s := fmt.Sprintf("%v:%v", t.PkgPath(), t.Name())
			s = strings.Replace(s, "github.com/", "github:", 1)
			s = strings.ReplaceAll(s, "/", "_")
			return s
		},
		Mapper: func(t reflect.Type) *JS.Schema {
			return nil
		},
		KeyNamer: strcase.LowerCamelCase,
	}
}

func DefaultAlterField(_ reflect.Type, f *reflect.StructField) bool {
	matchType := f.Type
	if matchType.Kind() == reflect.Ptr {
		matchType = f.Type.Elem()
	}
	switch matchType {
	case ToElemType((*net.IPOrDomain)(nil)):
		fallthrough
	case ToElemType((*cfgcommon.Address)(nil)):
		fallthrough
	case ToElemType((*cfgcommon.PortList)(nil)):
		fallthrough
	case ToElemType((*duration.Duration)(nil)):
		f.Type = ToElemType((*CustomString)(nil))

	case ToElemType((*cfgcommon.StringList)(nil)):
		fallthrough
	case ToElemType((*cfgcommon.NetworkList)(nil)):
		f.Type = ToElemType((*CustomStringList)(nil))

	case ToElemType((*cfgcommon.PortRange)(nil)):
		f.Type = ToElemType((*CustomPortRange)(nil))

	case ToElemType((*anypb.Any)(nil)):
		f.Type = ToElemType((*CustomAny)(nil))
	}
	return false
}

func LoadTypeByAlias(interfaceType, name string) reflect.Type {
	fsdef := envimpl.NewDefaultFileSystemDefaultImpl()
	ctx := envctx.ContextWithEnvironment(context.TODO(), fsdef)
	msg, err := registry.LoadImplementationByAlias(ctx, interfaceType, name, []byte("{}"))
	if err != nil {
		log.Fatalln(err)
	}
	return reflect.TypeOf(msg).Elem()
}

func ToElemType(ptr interface{}) reflect.Type {
	return reflect.TypeOf(ptr).Elem()
}

func BuildBasicObjectSchema(r *JS.Reflector, d JS.Definitions, t reflect.Type, excludes []string) *JS.Schema {
	res := r.RefOrReflectTypeToSchema(d, t)
	s := res
	if s.Ref != "" {
		defName := strings.Replace(s.Ref, "#/$defs/", "", 1)
		s = d[defName]
	}
	for _, name := range excludes {
		s.Properties.Delete(name)
	}
	return res
}

func BuildOneOfItemSchema(r *JS.Reflector, d JS.Definitions, idKey string, configKey string, idName string, nodeType reflect.Type) *JS.Schema {
	props := orderedmap.New()
	props.Set(idKey, &JS.Schema{Const: idName})
	if nodeType != nil {
		props.Set(configKey, r.RefOrReflectTypeToSchema(d, nodeType))
	}
	return &JS.Schema{
		Type:       "object",
		Properties: props,
	}
}

func BuildOneOfConfigsSchema(r *JS.Reflector, d JS.Definitions, idKey string, configKey string, interfaceType string, shortNames []string) *JS.Schema {
	var schemas []*JS.Schema
	for _, name := range shortNames {
		s := BuildOneOfItemSchema(r, d, idKey, configKey, name, LoadTypeByAlias(interfaceType, name))
		schemas = append(schemas, s)
	}
	return &JS.Schema{
		OneOf: schemas,
	}
}

func BuildSingleOrArraySchema(r *JS.Reflector, d JS.Definitions, t reflect.Type) *JS.Schema {
	s := r.RefOrReflectTypeToSchema(d, t)
	return &JS.Schema{OneOf: []*JS.Schema{
		s,
		{Type: "array", Items: s},
	}}
}

func BuildRouterStrategySchema(r *JS.Reflector, d JS.Definitions, idKey string, configKey string) *JS.Schema {
	return BuildOneOfConfigsSchema(r, d, idKey, configKey, "balancer", []string{
		"random", "leastping", "leastload",
	})
}
