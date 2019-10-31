package value

import (
	"reflect"
)

func (v *Value) getrv() reflect.Value {
	if v.rv.Kind() == reflect.Invalid {
		v.rv = reflect.ValueOf(v.iv)
	}

	return v.rv
}

func (v *Value) getiv() interface{} {
	if v.iv == nil {
		v.iv = v.rv.Interface()
	}
	return v.iv
}

//// get op

func (v *Value) mapGet(key interface{}) *Value {
	val := v.getrv().MapIndex(reflect.ValueOf(key))
	if val.Kind() == reflect.Invalid {
		return nil
	}
	return &Value{rv: val}
}

func (v *Value) sliceGet(idx int) *Value {
	rv := v.getrv()
	if idx >= rv.Len() {
		return nil
	}

	val := rv.Index(idx)
	return &Value{rv: val}
}

func (v *Value) structGet(fieldName string) *Value {
	field := v.getrv().FieldByName(fieldName)
	if field.Kind() == reflect.Invalid {
		return nil
	}
	return &Value{rv: field}
}

//// put op

func (v *Value) mapPut(key, val interface{}) error {
	if v.getrv().IsNil() {
		v.rv = reflect.MakeMap(v.getrv().Type())
		v.iv = v.rv.Interface()
	}
	v.getrv().SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
	return nil
}

func (v *Value) arrayPut(idx int, elem interface{}) error {
	rv := v.getrv()
	if idx >= rv.Cap() {
		return &ErrOutOfRange{"Value.arrayPut"}
	}
	ev := rv.Index(idx)
	ev.Set(reflect.ValueOf(elem))
	return nil
}

func (v *Value) slicePut(idx int, elem interface{}) error {
	rv := v.getrv()
	if idx < rv.Len() { // set
		ev := rv.Index(idx)
		ev.Set(reflect.ValueOf(elem))
	} else { // append
		x := reflect.ValueOf(elem)
		zv := reflect.Zero(x.Type())
		for i := rv.Len(); i < idx; i++ {
			rv = reflect.Append(rv, zv)
		}
		rv = reflect.Append(rv, x)

		v.rv = rv
		v.iv = rv.Interface()
	}
	return nil
}

// structPut ...
func (v *Value) structPut(fn string, field interface{}) error {
	fv := v.getrv().FieldByName(fn)
	if !fv.IsValid() {
		return &ErrNotExist{"Value.structPut", fn + " field"}
	}
	fv.Set(reflect.ValueOf(field))
	return nil
}

func (v *Value) bool() bool {
	return v.getrv().Bool()
}

func (v *Value) int() int64 {
	return v.getrv().Int()
}

func (v *Value) uint() uint64 {
	return v.getrv().Uint()
}

func (v *Value) float() float64 {
	return v.getrv().Float()
}

func (v *Value) complex_() complex128 {
	return v.getrv().Complex()
}

func (v *Value) string() string {
	return v.getrv().String()
}

func (v *Value) interface_() interface{} {
	return v.getrv().Interface()
}

//// map op

func (v *Value) mapMap() map[*Value]*Value {
	l := v.getrv().Len()
	m := make(map[*Value]*Value, l)
	iter := v.getrv().MapRange()
	for iter.Next() {
		key := iter.Key()
		val := iter.Value()
		m[&Value{rv: key}] = &Value{rv: val}
	}
	return m
}

func (v *Value) sliceMap() map[*Value]*Value {
	rv := v.getrv()
	m := make(map[*Value]*Value, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ev := rv.Index(i)
		m[&Value{iv: i}] = &Value{rv: ev}
	}
	return m
}

func (v *Value) structMap() map[*Value]*Value {
	num := v.getrv().NumField()
	m := make(map[*Value]*Value, num)
	rv := v.getrv()
	rt := rv.Type()
	for i := 0; i < num; i++ {
		fn := rt.Field(i).Name
		fv := rv.Field(i)
		m[&Value{iv: fn}] = &Value{rv: fv}
	}
	return m
}

//// slice op

func (v *Value) sliceSlice() []*Value {
	rv := v.getrv()
	s := make([]*Value, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		ev := rv.Index(i)
		s[i] = &Value{rv: ev}
	}
	return s
}

func (v *Value) structSlice() []*Value {
	rv := v.getrv()
	num := rv.NumField()
	s := make([]*Value, num)
	for i := 0; i < num; i++ {
		fv := rv.Field(i)
		s[i] = &Value{rv: fv}
	}
	return s
}

//// alist op

func (v *Value) mapAList() [][2]*Value {
	alist := make([][2]*Value, 0, v.getrv().Len())
	iter := v.getrv().MapRange()
	for iter.Next() {
		key := iter.Key()
		val := iter.Value()
		alist = append(alist, [2]*Value{&Value{rv: key}, &Value{rv: val}})
	}
	return alist
}

func (v *Value) sliceAList() [][2]*Value {
	l := v.getrv().Len()
	alist := make([][2]*Value, 0, l)
	rv := v.getrv()
	for i := 0; i < l; i++ {
		ev := rv.Index(i)
		alist = append(alist, [2]*Value{&Value{iv: i}, &Value{rv: ev}})
	}
	return alist
}

func (v *Value) structAList() [][2]*Value {
	num := v.getrv().NumField()
	alist := make([][2]*Value, 0, num)
	rv := v.getrv()
	rt := rv.Type()
	for i := 0; i < num; i++ {
		fn := rt.Field(i).Name
		fv := rv.Field(i)
		alist = append(alist, [2]*Value{&Value{iv: fn}, &Value{rv: fv}})
	}
	return alist
}

//// plist op

func (v *Value) mapPList() []*Value {
	l := v.getrv().Len()
	plist := make([]*Value, 0, 2*l)
	iter := v.getrv().MapRange()
	for iter.Next() {
		key := iter.Key()
		val := iter.Value()
		plist = append(plist, &Value{rv: key}, &Value{rv: val})
	}
	return plist
}

func (v *Value) slicePList() []*Value {
	l := v.getrv().Len()
	plist := make([]*Value, 0, 2*l)
	rv := v.getrv()
	for i := 0; i < l; i++ {
		ev := rv.Index(i)
		plist = append(plist, &Value{iv: i}, &Value{rv: ev})
	}
	return plist
}

func (v *Value) structPList() []*Value {
	num := v.getrv().NumField()
	plist := make([]*Value, 0, 2*num)
	rv := v.getrv()
	rt := rv.Type()
	for i := 0; i < num; i++ {
		fn := rt.Field(i).Name
		fv := rv.Field(i)
		plist = append(plist, &Value{iv: fn}, &Value{rv: fv})
	}
	return plist
}
