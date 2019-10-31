// Package value manipulates(gets,sets,convert etc.) interface{} value,
// convenient and simple, and reflect-based.
package value

import (
	"fmt"
	"math/bits"
	"reflect"
)

// Value ...
type Value struct {
	iv interface{}
	rv reflect.Value
}

// New new a Value from v
func New(v interface{}) *Value {
	return &Value{iv: v}
}

// Get returns the value with the given key.
//
// If t's kind is Map, Get returns the value associated with key in the map.
// If t's kind is Array or Slice, Get returns t's k'th element, the k must be int.
// If t's kind is Struct, Get returns the struct field with the given field name, the k must be string.
// if t's kind is Interface or Ptr, indirect it.
// It returns the nil if k is not found in the t.
// It returns error if t's kind is not Map, Array, Slice or Struct.
func (v *Value) Get(k interface{}) (*Value, error) {
	rv := v.getrv()
	switch rv.Kind() {
	case reflect.Map:
		return v.mapGet(k), nil
	case reflect.Array, reflect.Slice:
		return v.sliceGet(k.(int)), nil
	case reflect.Struct:
		return v.structGet(k.(string)), nil
	case reflect.Interface, reflect.Ptr:
		vv := &Value{rv: indirect(rv)}
		return vv.Get(k)
	default:
		return nil, &ErrUnsupportedKind{"Value.Get", rv.Kind()}
	}
}

// Set set t's value to v.
//
// If t's value can't setable, returns ErrCannotSet.
// If t's kind and v's kind is not equivalence, returns ErrTypeUnequal.
// It returns nil, that set successful.
//
// If set map key, struct field or array/slice index, using Value.Put.
func (v *Value) Set(iv interface{}) error {
	rv := v.getrv()
	if rv.Kind() == reflect.Interface || rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if !rv.CanSet() {
		return &ErrCannotSet{"Value.Set"}
	}

	irv := reflect.ValueOf(iv)
	if rv.Kind() != irv.Kind() {
		return &ErrTypeUnequal{"Value.Set", rv.Kind(), irv.Kind()}
	}

	rv.Set(irv)

	// reset
	v.iv = nil

	return nil
}

// Put put k, v to map, array, slice or struct(structed type).
//
// If t's kind is map, the k indicates key of map.
// If t's kind is array/slice, the k indecates index of array/slice.
// If t's kind is struct, the k indecates fieldname of struct.
//
// If k in t, and set k's value to v.
//
// If t's kind is not map, array, slice or struct, returns ErrUnsupportedKind.
func (v *Value) Put(key, val interface{}) (err error) {
	rv := v.getrv()

	switch rv.Kind() {
	case reflect.Map:
		return v.mapPut(key, val)
	case reflect.Slice:
		return v.slicePut(key.(int), val)
	case reflect.Ptr:
		rvv := indirect(rv)
		switch rvv.Kind() {
		case reflect.Array:
			return (&Value{rv: rvv}).arrayPut(key.(int), val)
		case reflect.Struct:
			return (&Value{rv: rvv}).structPut(key.(string), val)
		default:
			return &ErrUnsupportedKind{"Value.Put", v.getrv().Kind()}
		}
	default:
		return &ErrUnsupportedKind{"Value.Put", v.getrv().Kind()}
	}
}

// Bytes returns t's underlying value as a []bytes.
// It returns error if t's underlying value is not a slice of bytes.
func (v *Value) Bytes() ([]byte, error) {
	tv := v.getrv()
	switch tv.Kind() {
	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(tv)}).Bytes()
	case reflect.Slice:
		elemk := tv.Type().Elem().Kind()
		if elemk != reflect.Uint8 {
			return nil, &ErrUnsupportedKind{"Value.Bytes", "slice of " + elemk.String()}
		}
		return tv.Bytes(), nil
	default:
		return nil, &ErrUnsupportedKind{"Value.Bytes", v.getrv().Kind()}
	}
}

// Bool returns t's underlying value.
// It returns error if t's kind is not Bool.
func (v *Value) Bool() (bool, error) {
	switch v.getrv().Kind() {
	case reflect.Bool:
		return v.bool(), nil
	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).Bool()
	default:
		return false, &ErrUnsupportedKind{"Value.Int", v.getrv().Kind()}
	}
}

// Int returns t's underlying value as an int.
// It returns error if t's kind is not Int, Int8, Int16, Int32, Uint8 or Uint16,
// and if t's kind is Int64 or Uint32 also Int is 32 bits.
func (v *Value) Int() (i int, err error) {
	switch v.getrv().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		i = int(v.int())
	case reflect.Int64:
		if bits.UintSize == 64 { // if int size is 64
			i = int(v.int())
		} else { // if int size is 32
			err = &ErrUnsupportedKind{"Value.Int", v.getrv().Kind()}
		}

	case reflect.Uint8, reflect.Uint16:
		i = int(v.uint())

	case reflect.Uint32:
		if bits.UintSize == 64 {
			i = int(v.uint())
		} else {
			err = &ErrUnsupportedKind{"Value.Int", v.getrv().Kind()}
		}

	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).Int()

	default:
		err = &ErrUnsupportedKind{"Value.Int", v.getrv().Kind()}
	}
	return
}

// Int8 returns t's underlying value as an int8.
// It returns error if t's kind is not Int8.
func (v *Value) Int8() (int8, error) {
	switch v.getrv().Kind() {
	case reflect.Int8:
		return int8(v.int()), nil
	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).Int8()
	default:
		return 0, &ErrUnsupportedKind{"Value.Int8", v.getrv().Kind()}
	}
}

// Int16 returns t's underlying value as an int16.
// It returns error if t's kind is not Int, Int8, Int16, or Uint8.
func (v *Value) Int16() (int16, error) {
	switch v.getrv().Kind() {
	case reflect.Int8, reflect.Int16:
		return int16(v.int()), nil
	case reflect.Uint8:
		return int16(v.uint()), nil
	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).Int16()
	default:
		return 0, &ErrUnsupportedKind{"Value.Int16", v.getrv().Kind()}
	}
}

// Int32 returns t's underlying value as an int32.
// It returns error if t's kind is not Int, Int8, Int16, Int32, Uint8 or Uint16,
// and if t's kind is Int also Int is 64 bits.
func (v *Value) Int32() (int32, error) {
	switch v.getrv().Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32:
		return int32(v.int()), nil
	case reflect.Int:
		if bits.UintSize == 32 { // 32
			return int32(v.int()), nil
		}
		return 0, &ErrUnsupportedKind{"Value.Int32", v.getrv().Kind()}

	case reflect.Uint8, reflect.Uint16:
		return int32(v.uint()), nil

	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).Int32()

	default:
		return 0, &ErrUnsupportedKind{"Value.Int32", v.getrv().Kind()}
	}
}

// Int64 returns t's underlying value as an int64.
// It returns error if t's kind is not Int, Int8, Int16, Int32, Uint8, Uint16, Uint32
// and if t's kind is Uint also Uint is 64 bits.
func (v *Value) Int64() (int64, error) {
	switch v.getrv().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.int(), nil

	case reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return int64(v.uint()), nil
	case reflect.Uint:
		if bits.UintSize == 32 { // 32
			return int64(v.uint()), nil
		}
		return 0, &ErrUnsupportedKind{"Value.Int64", v.getrv().Kind()}

	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).Int64()

	default:
		return 0, &ErrUnsupportedKind{"Value.Int64", v.getrv().Kind()}
	}
}

// Uint returns t's underlying value as an uint.
// It returns error if t's kind is not Uint, Uint8, Uint16 or Uint32,
// and if t's kind is Uint64 also Uint is 32 bits.
func (v *Value) Uint() (i uint, err error) {
	switch v.getrv().Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		i = uint(v.uint())
	case reflect.Uint64:
		if bits.UintSize == 64 { // 64
			i = uint(v.uint())
		} else { // 32
			err = &ErrUnsupportedKind{"Value.Uint", v.getrv().Kind()}
		}
	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).Uint()

	default:
		err = &ErrUnsupportedKind{"Value.Uint", v.getrv().Kind()}
	}
	return
}

// Uint8 returns t's underlying value as an uint8.
// It returns error if t's kind is not Uint8.
func (v *Value) Uint8() (uint8, error) {
	switch v.getrv().Kind() {
	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).Uint8()
	case reflect.Uint8:
		return uint8(v.uint()), nil
	default:
		return 0, &ErrUnsupportedKind{"Value.Uint8", v.getrv().Kind()}
	}
}

// Uint16 returns t's underlying value as an uint16.
// It returns error if t's kind is not Uint8 or Uint16.
func (v *Value) Uint16() (uint16, error) {
	switch v.getrv().Kind() {
	case reflect.Uint8, reflect.Uint16:
		return uint16(v.uint()), nil
	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).Uint16()
	default:
		return 0, &ErrUnsupportedKind{"Value.Uint16", v.getrv().Kind()}
	}
}

// Uint32 returns t's underlying value as an uint32.
// It returns error if t's kind is not Uint8, Uint16 or Uint32,
// and if t's kind is Uint also Uint is 64 bits.
func (v *Value) Uint32() (uint32, error) {
	switch v.getrv().Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return uint32(v.uint()), nil

	case reflect.Uint:
		if bits.UintSize == 32 { // 32
			return uint32(v.uint()), nil
		}
		return 0, &ErrUnsupportedKind{"Value.Uint32", v.getrv().Kind()}

	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).Uint32()

	default:
		return 0, &ErrUnsupportedKind{"Value.Uint32", v.getrv().Kind()}
	}
}

// Uint64 returns t's underlying value as an uint64.
// It returns error if t's kind is not Uint*.
func (v *Value) Uint64() (uint64, error) {
	switch v.getrv().Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.uint(), nil
	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).Uint64()
	default:
		return 0, &ErrUnsupportedKind{"Value.Uint64", v.getrv().Kind()}
	}
}

// Float32 returns t's underlying value as an float32.
// It returns error if t's kind is not Uint*, Int* or Float32.
func (v *Value) Float32() (float32, error) {
	switch v.getrv().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float32(v.int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float32(v.uint()), nil
	case reflect.Float32:
		return float32(v.float()), nil
	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).Float32()
	default:
		return 0, &ErrUnsupportedKind{"Value.Float32", v.getrv().Kind()}
	}
}

// Float64 returns t's underlying value as an float64.
// It returns error if t's kind is not Uint*, Int* or Float*.
func (v *Value) Float64() (float64, error) {
	switch v.getrv().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.uint()), nil
	case reflect.Float32, reflect.Float64:
		return v.float(), nil
	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).Float64()
	default:
		return 0, &ErrUnsupportedKind{"Value.Float64", v.getrv().Kind()}
	}
}

// Complex64 returns t's underlying value as an complex64.
// It returns error if t's kind is not Uint*, Int*, Float32 or Complex64.
func (v *Value) Complex64() (complex64, error) {
	switch v.getrv().Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return complex(float32(v.uint()), 0), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return complex(float32(v.int()), 0), nil
	case reflect.Float32:
		return complex(float32(v.float()), 0), nil
	case reflect.Complex64:
		return complex64(v.complex_()), nil
	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).Complex64()
	default:
		return 0i, &ErrUnsupportedKind{"Value.Complex64", v.getrv().Kind()}
	}
}

// Complex128 returns t's underlying value as an complex128.
// It returns error if t's kind is not Uint*, Int*, Float* or Complex*.
func (v *Value) Complex128() (complex128, error) {
	switch v.getrv().Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return complex(float64(v.uint()), 0), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return complex(float64(v.int()), 0), nil
	case reflect.Float32, reflect.Float64:
		return complex(v.float(), 0), nil
	case reflect.Complex64, reflect.Complex128:
		return v.complex_(), nil
	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).Complex128()
	default:
		return 0i, &ErrUnsupportedKind{"Value.Complex128", v.getrv().Kind()}
	}
}

// Map returns t's underlying value as a map.
// It returns error if t's kind is not Map, Array, Slice or Struct.
func (v *Value) Map() (map[*Value]*Value, error) {
	switch v.getrv().Kind() {
	case reflect.Map:
		return v.mapMap(), nil
	case reflect.Array, reflect.Slice:
		return v.sliceMap(), nil
	case reflect.Struct:
		return v.structMap(), nil
	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).Map()
	default:
		return nil, &ErrUnsupportedKind{"Value.Map", v.getrv().Kind()}
	}
}

// Slice returns t's underlying value as a slice.
// It returns error if t's kind is not Array, Slice or Struct.
func (v *Value) Slice() ([]*Value, error) {
	switch v.getrv().Kind() {
	case reflect.Array, reflect.Slice:
		return v.sliceSlice(), nil
	case reflect.Struct:
		return v.structSlice(), nil
	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).Slice()
	default:
		return nil, &ErrUnsupportedKind{"Value.Slice", v.getrv().Kind()}
	}
}

// AList returns t's underlying value as an association list.
// It returns error if t's kind is not Map, Array, Slice or Struct.
func (v *Value) AList() ([][2]*Value, error) {
	switch v.getrv().Kind() {
	case reflect.Map:
		return v.mapAList(), nil
	case reflect.Array, reflect.Slice:
		return v.sliceAList(), nil
	case reflect.Struct:
		return v.structAList(), nil
	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).AList()
	default:
		return nil, &ErrUnsupportedKind{"Value.AList", v.getrv().Kind()}
	}
}

// PList returns t's underlying value as an property list.
// It returns error if t's kind is not Map, Array, Slice or Struct.
func (v *Value) PList() ([]*Value, error) {
	switch v.getrv().Kind() {
	case reflect.Map:
		return v.mapPList(), nil
	case reflect.Array, reflect.Slice:
		return v.slicePList(), nil
	case reflect.Struct:
		return v.structPList(), nil
	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).PList()
	default:
		return nil, &ErrUnsupportedKind{"Value.PList", v.getrv().Kind()}
	}
}

func (v *Value) Interface() interface{} {
	return v.getiv()
}

func (v *Value) Ptr() uintptr {
	return v.getrv().Pointer()
}

var _StringerType = reflect.TypeOf((fmt.Stringer)(nil))

func (v *Value) String() (string, error) {
	// TODO: check if implement string

	switch v.getrv().Kind() {
	case reflect.Invalid:
		return "", nil
	case reflect.String:
		return v.string(), nil
	case reflect.Bool:
		return fmt.Sprintf("%t", v.bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return fmt.Sprintf("%d", v.uint()), nil
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%g", v.float()), nil
	case reflect.Complex64, reflect.Complex128:
		return fmt.Sprintf("%g", v.complex_()), nil
	case reflect.Chan:
		return fmt.Sprintf("%p", v.getrv().Interface()), nil
	case reflect.UnsafePointer:
		return fmt.Sprintf("%x", v.getrv().Interface()), nil
	case reflect.Struct:
		return fmt.Sprintf("%#v", v.getrv().Interface()), nil
	case reflect.Slice, reflect.Array, reflect.Map:
		return fmt.Sprintf("%v", v.getrv().Interface()), nil
	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).String()
	default:
		return "", &ErrUnsupportedKind{"Value.String", v.getrv().Kind()}
	}
}

type eachDoFunc func(k, v *Value) error

// EachDo ...
func (v *Value) EachDo(f eachDoFunc) error {
	switch v.getrv().Kind() {
	case reflect.Map, reflect.Array, reflect.Slice, reflect.Struct:
		m, err := v.Map()
		if err != nil {
			return err
		}
		for k, v := range m {
			if err := f(k, v); err != nil {
				return err
			}
		}
	case reflect.Chan:
		idx := 0
		for {
			v, ok := v.getrv().Recv()
			if !ok {
				break
			}

			if err := f(&Value{iv: idx}, &Value{rv: v}); err != nil {
				return err
			}

			idx++
		}
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:

		if err := f(nil, v); err != nil {
			return err
		}

	case reflect.Interface, reflect.Ptr:
		return (&Value{rv: indirect(v.getrv())}).EachDo(f)

	default:
		return &ErrUnsupportedKind{"Value.EachDo", v.getrv().Kind()}
	}
	return nil
}
