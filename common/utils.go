package common

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	jss "github.com/invopop/jsonschema"
	"github.com/stoewer/go-strcase"

	"github.com/v2fly/v2ray-core/v5/common/environment/envctx"
	"github.com/v2fly/v2ray-core/v5/common/environment/envimpl"
	"github.com/v2fly/v2ray-core/v5/common/registry"
)

func NewDefaultReflector() jss.Reflector {
	return jss.Reflector{
		RequiredFromJSONSchemaTags: true,
		AllowAdditionalProperties:  true,
		Namer: func(t reflect.Type) string {
			s := fmt.Sprintf("%v:%v", t.PkgPath(), t.Name())
			s = strings.Replace(s, "github.com/", "github:", 1)
			s = strings.ReplaceAll(s, "/", "_")
			return s
		},
		KeyNamer: strcase.LowerCamelCase,
	}
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
