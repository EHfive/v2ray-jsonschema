package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	strcase "github.com/stoewer/go-strcase"

	"github.com/EHfive/v2ray-jsonschema/v5config"
	"github.com/invopop/jsonschema"
)

func main() {
	reflector := &jsonschema.Reflector{
		DoNotReference:            false,
		AllowAdditionalProperties: true,
		Namer: func(t reflect.Type) string {
			s := fmt.Sprintf("%v:%v", t.PkgPath(), t.Name())
			s = strings.Replace(s, "github.com/", "github:", 1)
			s = strings.ReplaceAll(s, "/", "_")
			return s
		},
		Mapper: func(t reflect.Type) *jsonschema.Schema {
			// log.Println(t.Name())
			return nil
		},
		KeyNamer: strcase.LowerCamelCase,
	}
	//schema := reflector.Reflect(&v5config.V5Config{})

	schema := v5config.JSONSchema(*reflector)
	data, _ := json.MarshalIndent(schema, "", "  ")
	fmt.Println(string(data))
}
