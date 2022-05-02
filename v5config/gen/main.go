package main

import (
	"encoding/json"
	"fmt"

	"github.com/EHfive/v2ray-jsonschema/common"
	"github.com/EHfive/v2ray-jsonschema/v5config"
)

func main() {
	r := common.NewDefaultReflector()
	// r.DoNotReference = true
	schema := v5config.JSONSchema(r)
	data, _ := json.MarshalIndent(schema, "", "  ")
	fmt.Println(string(data))
}
