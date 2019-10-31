package value

import (
	"reflect"
)

func indirect(v reflect.Value) reflect.Value {
	for {
		if v.Kind() != reflect.Interface && v.Kind() != reflect.Ptr {
			return v
		}
		v = v.Elem()
	}
}
