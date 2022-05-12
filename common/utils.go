package common

import (
	"context"
	"encoding/json"
	"log"
	"reflect"

	"github.com/v2fly/v2ray-core/v5/common/environment/envctx"
	"github.com/v2fly/v2ray-core/v5/common/environment/envimpl"
	"github.com/v2fly/v2ray-core/v5/common/registry"
)

func ToElemType(ptr interface{}) reflect.Type {
	return reflect.TypeOf(ptr).Elem()
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

var rawMessageType = reflect.TypeOf(json.RawMessage{})

func ReplaceTypeElem(holder reflect.Type, oldElem reflect.Type, newElem reflect.Type) (reflect.Type, bool) {
	var r func(reflect.Type) reflect.Type
	replaced := false
	r = func(t reflect.Type) reflect.Type {
		if t == oldElem {
			replaced = true
			return newElem
		}
		switch t.Kind() {
		case reflect.Array:
			return reflect.ArrayOf(t.Len(), r(t.Elem()))
		case reflect.Slice:
			if t == rawMessageType {
				break
			}
			return reflect.SliceOf(r(t.Elem()))
		case reflect.Map:
			return reflect.MapOf(t.Key(), r(t.Elem()))
		case reflect.Chan:
			return reflect.ChanOf(t.ChanDir(), r(t.Elem()))
		case reflect.Ptr:
			return r(t.Elem())
		}
		if oldElem == nil {
			replaced = true
			return newElem
		}
		return t
	}
	return r(holder), replaced
}

type ReplaceTypePair [2]interface{}

func ReplaceTypeElemByPairs(holder reflect.Type, pairs []ReplaceTypePair) (reflect.Type, bool) {
	for _, pair := range pairs {
		t, ok := ReplaceTypeElem(holder, ToElemType(pair[0]), ToElemType(pair[1]))
		if ok {
			return t, true
		}
	}
	return holder, false
}

type ReplaceFieldTypePair [3]interface{}

func ReplaceFieldTypeElemByPairs(holder reflect.Type, field reflect.StructField, pairs []ReplaceFieldTypePair) (reflect.Type, bool) {
	for _, pair := range pairs {
		if pair[1].(string) == field.Name && holder == ToElemType(pair[0]) {
			return ReplaceTypeElem(field.Type, nil, ToElemType(pair[2]))
		}
	}
	return field.Type, false
}
