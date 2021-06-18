package value

// MustInt must api for Int()
func (v *Value) MustInt() int {
	i, err := v.Int()
	if err != nil {
		panic(err)
	}
	return i
}

// MustInt8 must api for Int8()
func (v *Value) MustInt8() int8 {
	i, err := v.Int8()
	if err != nil {
		panic(err)
	}
	return i
}

// MustInt16 must api for Int16()
func (v *Value) MustInt16() int16 {
	i, err := v.Int16()
	if err != nil {
		panic(err)
	}
	return i
}

// MustInt32 must api for Int32()
func (v *Value) MustInt32() int32 {
	i, err := v.Int32()
	if err != nil {
		panic(err)
	}
	return i
}

// MustInt64 must api for Int64()
func (v *Value) MustInt64() int64 {
	i, err := v.Int64()
	if err != nil {
		panic(err)
	}
	return i
}

// MustUint must api for Uint()
func (v *Value) MustUint() uint {
	i, err := v.Uint()
	if err != nil {
		panic(err)
	}
	return i
}

// MustUint8 must api for Uint8()
func (v *Value) MustUint8() uint8 {
	i, err := v.Uint8()
	if err != nil {
		panic(err)
	}
	return i
}

// MustUint16 must api for Uint16()
func (v *Value) MustUint16() uint16 {
	i, err := v.Uint16()
	if err != nil {
		panic(err)
	}
	return i
}

// MustUint32 must api for Uint32()
func (v *Value) MustUint32() uint32 {
	i, err := v.Uint32()
	if err != nil {
		panic(err)
	}
	return i
}

// MustUint64 must api for Uint64()
func (v *Value) MustUint64() uint64 {
	i, err := v.Uint64()
	if err != nil {
		panic(err)
	}
	return i
}

// MustFloat32 must api for Float32()
func (v *Value) MustFloat32() float32 {
	val, err := v.Float32()
	if err != nil {
		panic(err)
	}
	return val
}

// MustFloat64 must api for Float64()
func (v *Value) MustFloat64() float64 {
	val, err := v.Float64()
	if err != nil {
		panic(err)
	}
	return val
}

// MustComplex64 must api for Complex64()
func (v *Value) MustComplex64() complex64 {
	val, err := v.Complex64()
	if err != nil {
		panic(err)
	}
	return val
}

// MustComplex128 must api for Complex128()
func (v *Value) MustComplex128() complex128 {
	val, err := v.Complex128()
	if err != nil {
		panic(err)
	}
	return val
}

// MustBool must api for Bool()
func (v *Value) MustBool() bool {
	val, err := v.Bool()
	if err != nil {
		panic(err)
	}
	return val
}

// MustBytes must api for Bytes()
func (v *Value) MustBytes() []byte {
	val, err := v.Bytes()
	if err != nil {
		panic(err)
	}
	return val
}

// MustString must api for String()
func (v *Value) MustString() string {
	val, err := v.String()
	if err != nil {
		panic(err)
	}
	return val
}

// MustGet must api for Get
func (v *Value) MustGet(k interface{}) *Value {
	val, err := v.Get(k)
	if err != nil {
		panic(err)
	}
	return val
}

// MustMap must api for Map
func (v *Value) MustMap() map[*Value]*Value {
	tm, err := v.Map()
	if err != nil {
		panic(err)
	}
	return tm
}

// MustSlice must api for Slice
func (v *Value) MustSlice() []*Value {
	ts, err := v.Slice()
	if err != nil {
		panic(err)
	}
	return ts
}

// MustAList must api for AList
func (v *Value) MustAList() [][2]*Value {
	tl, err := v.AList()
	if err != nil {
		panic(err)
	}
	return tl
}

// MustPList must api for PList
func (v *Value) MustPList() []*Value {
	tl, err := v.PList()
	if err != nil {
		panic(err)
	}
	return tl
}
