package main

import (
	"encoding/json"
	"fmt"

	"github.com/EHfive/v2ray-jsonschema/common"
	"github.com/EHfive/v2ray-jsonschema/v4config"
)

func main() {
	r := common.NewDefaultReflector()
	s := v4config.JSONSchema(r)
	data, _ := json.MarshalIndent(s, "", "  ")
	fmt.Println(string(data))
}
