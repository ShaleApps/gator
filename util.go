package gator

import (
	"errors"
	"reflect"
)

func isStructOrStructPtr(src interface{}) error {
	objT := reflect.TypeOf(src)
	objV := reflect.ValueOf(src)
	switch {
	case isStruct(objT), isStructPtr(objT):
		objT = objT.Elem()
		objV = objV.Elem()
	default:
		return errors.New("gator: src must be a struct or a pointer to a struct")
	}
	return nil
}

func isStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct
}

func isStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

func isArrayOrSlice(a interface{}) bool {
	if a == nil {
		return false
	}
	switch reflect.TypeOf(a).Kind() {
	case reflect.Array, reflect.Slice:
		return true
	default:
		return false
	}
}

func lengthOf(a interface{}) (int, bool) {
	if a == nil {
		return 0, false
	}
	switch reflect.TypeOf(a).Kind() {
	case reflect.Map, reflect.Array, reflect.String, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(a).Len(), true
	default:
		return 0, false
	}
}
