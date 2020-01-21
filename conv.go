package value

import (
	"fmt"
	"math/bits"
	"net"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var (
	// TimeLayout default time layout
	TimeLayout = "Mon Jan 2 15:04:05 -0700 MST 2006"
)

var (
	intLevel = map[reflect.Kind]int{
		reflect.Int8:  1,
		reflect.Int16: 2,
		reflect.Int32: 3,
		reflect.Int64: 4,
	}
	uintLevel = map[reflect.Kind]int{
		reflect.Uint8:  1,
		reflect.Uint16: 2,
		reflect.Uint32: 3,
		reflect.Uint64: 4,
	}
	floatLevel = map[reflect.Kind]int{
		reflect.Float32: 1,
		reflect.Float64: 2,
	}
	complexLevel = map[reflect.Kind]int{
		reflect.Complex64:  1,
		reflect.Complex128: 2,
	}
)

// ConvTo convToert t to dst
func (v *Value) ConvTo(dst interface{}) error {
	dstv := reflect.ValueOf(dst)
	if dstv.Kind() != reflect.Ptr {
		return &ErrUnsupportedKind{"Value.ConvTo", dstv.Kind()}
	}
	return v.convTo(dstv.Elem())
}

func (v *Value) convTo(dst reflect.Value) (err error) {
	switch dst.Type().String() {
	case "time.Duration":
		return v.convToTimeDuration(dst)
	case "time.Time":
		return v.convToTimeTime(dst)
	case "net.IP":
		return v.convToNetIP(dst)
	case "url.URL":
		return v.convToNetURL(dst)
	case "mail.Address":
		return v.convToMailAddress(dst)
	case "regexp.Regexp":
		return v.convToRegexpRegexp(dst)
	}

	switch dst.Kind() {
	case reflect.Bool:
		return v.convToBool(dst)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.convToInt(dst)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.convToUint(dst)

	case reflect.Float32, reflect.Float64:
		return v.convToFloat(dst)

	case reflect.Complex64, reflect.Complex128:
		return v.convToComplex(dst)

	case reflect.String:
		return v.convToString(dst)

	case reflect.Map:
		return v.convToMap(dst)

	case reflect.Array:
		return v.convToArray(dst)

	case reflect.Slice:
		return v.convToSlice(dst)

	case reflect.Struct:
		return v.convToStruct(dst)

	case reflect.Interface:
		return v.convToInterface(dst)

	case reflect.Ptr:
		return v.convToPtr(dst)

	default:
		return &ErrUnsupportedKind{"Value.convTo", dst.Kind()}
	}
}

func (v *Value) convToTimeDuration(dst reflect.Value) error {
	s, err := v.String()
	if err != nil {
		return err
	}
	td, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	dst.SetInt(int64(td))
	return nil
}

func (v *Value) convToTimeTime(dst reflect.Value) error {
	s, err := v.String()
	if err != nil {
		return err
	}
	t, err := time.Parse(TimeLayout, s)
	if err != nil {
		return err
	}
	dst.Set(reflect.ValueOf(t))
	return nil
}

func (v *Value) convToNetIP(dst reflect.Value) error {
	s, err := v.String()
	if err != nil {
		return err
	}
	ip := net.ParseIP(s)
	if ip == nil {
		return fmt.Errorf("invalid ip")
	}
	dst.Set(reflect.ValueOf(ip))
	return nil
}

func (v *Value) convToNetURL(dst reflect.Value) error {
	s, err := v.String()
	if err != nil {
		return err
	}
	url, err := url.Parse(s)
	if err != nil {
		return err
	}
	dst.Set(reflect.ValueOf(*url))
	return nil
}

func (v *Value) convToMailAddress(dst reflect.Value) error {
	s, err := v.String()
	if err != nil {
		return err
	}
	addr, err := mail.ParseAddress(s)
	if err != nil {
		return err
	}
	dst.Set(reflect.ValueOf(*addr))
	return nil
}

func (v *Value) convToRegexpRegexp(dst reflect.Value) error {
	s, err := v.String()
	if err != nil {
		return err
	}

	r, err := regexp.Compile(s)
	if err == nil {
		dst.Set(reflect.ValueOf(*r))
		return nil
	}

	r, err = regexp.CompilePOSIX(s)
	if err == nil {
		dst.Set(reflect.ValueOf(*r))
		return nil
	}

	return err
}

func (v *Value) convToPtr(dst reflect.Value) error {
	realdst := dst
	if dst.IsNil() {
		realdst = reflect.New(dst.Type().Elem())
	}
	if err := v.convTo(realdst.Elem()); err != nil {
		return err
	}
	dst.Set(realdst)

	return nil
}

func (v *Value) convToInterface(dst reflect.Value) error {
	dst.Set(v.getrv())
	return nil
}

func (v *Value) convToBool(dst reflect.Value) error {
	b, err := v.Bool()
	if err != nil {
		return err
	}

	dst.SetBool(b)
	return nil
}

func (v *Value) convToInt(dst reflect.Value) error {
	vk := v.getrv().Kind()
	dstk := dst.Kind()

	if dstk == reflect.Int { // convToert Int to Int32 or Int64
		fixedk := reflect.Invalid // fixed size kind
		if bits.UintSize == 32 {
			fixedk = reflect.Int32
		}
		if bits.UintSize == 64 {
			fixedk = reflect.Int64
		}
		dstk = fixedk
	}

	if intLevel[dstk] < intLevel[vk] {
		return &ErrTypeUnequal{
			"Value.convToInt",
			dstk,
			vk,
		}
	}

	iv, err := v.Int64()
	if err != nil {
		return err
	}
	dst.SetInt(iv)
	return nil
}

func (v *Value) convToUint(dst reflect.Value) error {
	vk := v.getrv().Kind()
	dstk := dst.Kind()

	if dstk == reflect.Uint { // convToert Uint to Uint32 or Uint64
		fixedk := reflect.Invalid
		if bits.UintSize == 32 {
			fixedk = reflect.Uint32
		}
		if bits.UintSize == 64 {
			fixedk = reflect.Uint64
		}
		dstk = fixedk
	}

	if uintLevel[dstk] < uintLevel[vk] {
		return &ErrTypeUnequal{
			"Value.convToUint",
			dstk,
			vk,
		}
	}

	uv, err := v.Uint64()
	if err != nil {
		return err
	}
	dst.SetUint(uv)
	return nil
}

func (v *Value) convToFloat(dst reflect.Value) error {
	vk := v.getrv().Kind()
	dstk := dst.Kind()

	if floatLevel[dstk] < floatLevel[vk] {
		return &ErrTypeUnequal{
			"Value.convToFloat",
			dstk,
			vk,
		}
	}

	fv, err := v.Float64()
	if err != nil {
		return err
	}
	dst.SetFloat(fv)
	return nil
}

func (v *Value) convToComplex(dst reflect.Value) error {
	vk := v.getrv().Kind()
	dstk := dst.Kind()

	if complexLevel[dstk] < complexLevel[vk] {
		return &ErrTypeUnequal{
			"Value.convToComplex",
			vk,
			dstk,
		}
	}

	cv, err := v.Complex128()
	if err != nil {
		return err
	}
	dst.SetComplex(cv)
	return nil
}

func (v *Value) convToString(dst reflect.Value) error {
	s, err := v.String()
	if err != nil {
		return err
	}
	dst.SetString(s)
	return nil
}

func (v *Value) convToMap(dst reflect.Value) error {
	if dst.IsNil() {
		dst.Set(reflect.MakeMap(dst.Type()))
	}

	vm, err := v.Map()
	if err != nil {
		return err
	}
	for srck, srcv := range vm {
		dstk := srck.getrv()
		dstv := dst.MapIndex(dstk)
		if dstv.Kind() == reflect.Invalid {
			dstv = reflect.New(dst.Type().Elem())
			dstv = dstv.Elem()
		}

		if err := srcv.convTo(dstv); err != nil {
			return err
		}
		dst.SetMapIndex(dstk, dstv)
	}
	return nil
}

func (v *Value) convToArray(dst reflect.Value) error {
	vs, err := v.Slice()
	if err != nil {
		return err
	}

	for idx, elem := range vs {
		if idx >= dst.Cap() {
			break
		}

		ev := dst.Index(idx)

		if err := elem.convTo(ev); err != nil {
			return err
		}
	}
	return nil
}

func (v *Value) convToSlice(dst reflect.Value) error {
	vs, err := v.Slice()
	if err != nil {
		return err
	}

	newSlice := dst
	for i, v := range vs {
		var ev reflect.Value
		if i < newSlice.Len() {
			ev = newSlice.Index(i)
		} else {
			ev = reflect.New(newSlice.Type().Elem())
			newSlice = reflect.Append(newSlice, ev.Elem())
			ev = newSlice.Index(i)
		}

		if err := v.convTo(ev); err != nil {
			return err
		}
	}
	dst.Set(newSlice)
	return nil
}
func (v *Value) convToStruct(dst reflect.Value) error {
	vm, err := v.Map()
	if err != nil {
		return err
	}

	dstype := dst.Type()
	tag2Fname := map[string]string{}
	lower2Fname := map[string]string{}
	passedFnames := map[string]bool{}
	for i := 0; i < dstype.NumField(); i++ {
		field := dstype.Field(i)
		tag := field.Tag.Get("value")
		if tag == "_" {
			passedFnames[field.Name] = true
		} else if tag != "" {
			tag2Fname[tag] = field.Name
		} else {
			lower2Fname[strings.ToLower(field.Name)] = field.Name
		}
	}

	matchCase := false
	for kv := range vm {
		key, err := kv.String()
		if err != nil {
			return err
		}
		if strings.ToLower(key) != key {
			matchCase = true
		}
	}
	for kv, vv := range vm {
		key, err := kv.String()
		if err != nil {
			return err
		}
		if passedFnames[key] {
			continue
		}
		fn := tag2Fname[key]
		if fn == "" {
			if matchCase {
				fn = key
			} else {
				fn = lower2Fname[key]
			}
		}
		f := dst.FieldByName(fn)
		if f.Kind() == reflect.Invalid {
			continue
		}
		if err := vv.convTo(f); err != nil {
			return err
		}
		passedFnames[fn] = true
	}
	return nil
}
