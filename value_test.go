package value

import (
	"math/bits"
	"net"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func ExpectErr(rets ...interface{}) Assertion {
	return Expect(rets[1])
}

var _ = Describe("Gets", func() {
	Context("with Bool()", func() {
		Specify("from bool kind", func() {
			b := true
			v := New(b)
			Expect(v.Bool()).To(Equal(b))
		})
		Specify("from bool ptr kind", func() {
			b := true
			v := New(&b)
			Expect(v.Bool()).To(Equal(b))
		})
		Specify("from other kind", func() {
			v := New("test")
			ExpectErr(v.Bool()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with Bytes()", func() {
		Specify("from []byte type", func() {
			x := []byte("abcd")
			v := New(x)
			Expect(v.Bytes()).To(Equal(x))
		})
		Specify("from []byte ptr type", func() {
			x := []byte("abcd")
			v := New(&x)
			Expect(v.Bytes()).To(Equal(x))
		})
		Specify("from other kind", func() {
			v := New("test")
			ExpectErr(v.Bytes()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))

			v = New([]int{1, 2, 3})
			ExpectErr(v.Bytes()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with Uint()", func() {
		Specify("from uint type", func() {
			x := uint(12)
			vs := []*Value{
				New(x),
				New(uint8(x)),
				New(uint16(x)),
				New(uint32(x)),
			}
			for _, v := range vs {
				Expect(v.Uint()).To(Equal(x))
			}

			ut64 := New(uint64(x))
			if bits.UintSize == 64 {
				Expect(ut64.Uint()).To(Equal(x))
			} else {
				ExpectErr(ut64.Uint()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
			}
		})
		Specify("from ptr kind", func() {
			x := uint(1)
			v := New(&x)
			Expect(v.Uint()).Should(Equal(x))
		})
		Specify("from other type", func() {
			v := New("test")
			ExpectErr(v.Uint()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with Uint8()", func() {
		Specify("from uint8 kind", func() {
			x := uint8(12)
			vs := []*Value{
				New(uint8(x)),
			}
			for _, v := range vs {
				Expect(v.Uint8()).To(Equal(x))
			}
		})
		Specify("from ptr kind", func() {
			x := uint8(1)
			v := New(&x)
			Expect(v.Uint8()).Should(Equal(x))
		})
		Specify("from other kind", func() {
			sv := New("test")
			ExpectErr(sv.Uint8()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with Uint16()", func() {
		Specify("from uint16 kind", func() {
			x := uint16(12)
			vs := []*Value{
				New(uint8(x)),
				New(uint16(x)),
			}
			for _, v := range vs {
				Expect(v.Uint16()).To(Equal(x))
			}
		})
		Specify("from ptr kind", func() {
			x := uint16(12)
			v := New(&x)
			Expect(v.Uint16()).Should(Equal(x))
		})
		Specify("from other kind", func() {
			sv := New("test")
			ExpectErr(sv.Uint16()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with Uint32()", func() {
		Specify("from uint32 kind", func() {
			x := uint32(12)
			vs := []*Value{
				New(uint8(x)),
				New(uint16(x)),
				New(uint32(x)),
			}
			for _, v := range vs {
				Expect(v.Uint32()).To(Equal(x))
			}

			uv := New(uint(x))
			if bits.UintSize == 32 {
				Expect(uv.Uint32()).To(Equal(x))
			} else {
				ExpectErr(uv.Uint32()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
			}
		})
		Specify("from ptr kind", func() {
			x := uint32(12)
			v := New(&x)
			Expect(v.Uint32()).Should(Equal(x))
		})
		Specify("from other kind", func() {
			sv := New("test")
			ExpectErr(sv.Uint32()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with Uint64()", func() {
		Specify("from uint* kind", func() {
			x := uint64(12)
			vs := []*Value{
				New(uint(x)),
				New(uint8(x)),
				New(uint16(x)),
				New(uint32(x)),
				New(uint64(x)),
			}
			for _, v := range vs {
				Expect(v.Uint64()).To(Equal(x))
			}
		})
		Specify("from ptr kind", func() {
			x := uint64(12)
			v := New(&x)
			Expect(v.Uint64()).Should(Equal(x))
		})
		Specify("from other kind", func() {
			sv := New("test")
			ExpectErr(sv.Uint64()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with Int()", func() {
		Specify("from int, int8, int16, int32, int64, uint8, uint16, uint32 kind", func() {
			x := 12
			vs := []*Value{
				New(int(x)),
				New(int8(x)),
				New(int16(x)),
				New(int32(x)),
				New(uint8(x)),
				New(uint16(x)),
			}
			for _, v := range vs {
				Expect(v.Int()).To(Equal(x))
			}

			t64 := New(int64(x))
			if bits.UintSize == 64 {
				Expect(t64.Int()).To(Equal(x))
			} else {
				ExpectErr(t64.Int()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
			}

			ut32 := New(uint32(x))
			if bits.UintSize == 64 {
				Expect(ut32.Int()).To(Equal(x))
			} else {
				ExpectErr(ut32.Int()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
			}
		})
		Specify("from ptr kind", func() {
			x := int(12)
			v := New(&x)
			Expect(v.Int()).Should(Equal(x))
		})
		Specify("from other kind", func() {
			v := New("test")
			ExpectErr(v.Int()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with Int8()", func() {
		Specify("from int8 kind", func() {
			t8 := New(int8(12))
			Expect(t8.Int8()).To(Equal(int8(12)))
		})
		Specify("from ptr kind", func() {
			x := int8(12)
			v := New(&x)
			Expect(v.Int8()).Should(Equal(x))
		})
		Specify("from other kind", func() {
			s := New("test")
			ExpectErr(s.Int8()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with Int16()", func() {
		Specify("from int8, int16, uint8 kind", func() {
			x := int16(12)
			vs := []*Value{
				New(int8(x)),
				New(int16(x)),
				New(uint8(x)),
			}
			for _, v := range vs {
				Expect(v.Int16()).To(Equal(x))
			}
		})
		Specify("from ptr kind", func() {
			x := int16(12)
			v := New(&x)
			Expect(v.Int16()).Should(Equal(x))
		})
		Specify("from other kind", func() {
			v := New("test")
			ExpectErr(v.Int16()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with Int32()", func() {
		Specify("from int, int{8,16,32}, uint{8, 16}", func() {
			x := int32(12)
			vs := []*Value{
				New(int8(x)),
				New(int16(x)),
				New(int32(x)),
				New(uint8(x)),
				New(uint16(x)),
			}
			for _, v := range vs {
				Expect(v.Int32()).To(Equal(x))
			}

			v := New(int(x))
			if bits.UintSize == 32 {
				Expect(v.Int32()).To(Equal(x))
			} else {
				ExpectErr(v.Int32()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
			}
		})
		Specify("from ptr kind", func() {
			x := int32(12)
			v := New(&x)
			Expect(v.Int32()).Should(Equal(x))
		})
		Specify("from other kind", func() {
			sv := New("test")
			ExpectErr(sv.Int32()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with Int64()", func() {
		Specify("from int, int{8,16,32,64}, uint, uint{8, 16, 32}", func() {
			x := int64(12)
			vs := []*Value{
				New(int(x)),
				New(int8(x)),
				New(int16(x)),
				New(int32(x)),
				New(int64(x)),
				New(uint8(x)),
				New(uint16(x)),
				New(uint32(x)),
			}
			for _, v := range vs {
				Expect(v.Int64()).To(Equal(x))
			}

			uv := New(uint(x))
			if bits.UintSize == 32 {
				Expect(uv.Int64()).To(Equal(x))
			} else {
				ExpectErr(uv.Int64()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
			}
		})
		Specify("from ptr kind", func() {
			x := int64(12)
			v := New(&x)
			Expect(v.Int64()).Should(Equal(x))
		})
		Specify("from other kind", func() {
			v := New("test")
			ExpectErr(v.Int64()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with Float32()", func() {
		Specify("from int*, uint*, float32 kind", func() {
			x := float32(12)
			vs := []*Value{
				New(x),
				New(int(x)),
				New(int8(x)),
				New(int16(x)),
				New(int32(x)),
				New(int64(x)),
				New(uint(x)),
				New(uint8(x)),
				New(uint16(x)),
				New(uint32(x)),
				New(uint64(x)),
			}
			for _, v := range vs {
				Expect(v.Float32()).To(Equal(x))
			}
		})
		Specify("from ptr kind", func() {
			x := float32(12)
			v := New(&x)
			Expect(v.Float32()).Should(Equal(x))
		})
		Specify("from other kind", func() {
			sv := New("test")
			ExpectErr(sv.Float32()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with Float64()", func() {
		Specify("from int*, uint*, float* kind", func() {
			x := float64(12)
			vs := []*Value{
				New(x),
				New(float32(x)),
				New(int(x)),
				New(int8(x)),
				New(int16(x)),
				New(int32(x)),
				New(int64(x)),
				New(uint(x)),
				New(uint8(x)),
				New(uint16(x)),
				New(uint32(x)),
				New(uint64(x)),
			}
			for _, v := range vs {
				Expect(v.Float64()).To(Equal(x))
			}
		})
		Specify("from ptr kind", func() {
			x := float64(12)
			v := New(&x)
			Expect(v.Float64()).Should(Equal(x))
		})
		Specify("from other kind", func() {
			sv := New("test")
			ExpectErr(sv.Float64()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with Complex64()", func() {
		Specify("from int*, uint*, float*, complex64 kind", func() {
			c := complex(float32(12), float32(13))
			r := 12
			vs := []*Value{
				New(int(r)),
				New(int8(r)),
				New(int16(r)),
				New(int32(r)),
				New(int64(r)),
				New(uint(r)),
				New(uint8(r)),
				New(uint16(r)),
				New(uint32(r)),
				New(uint64(r)),
				New(float32(r)),
			}
			for _, v := range vs {
				Expect(v.Complex64()).To(Equal(complex(float32(r), 0)))
			}
			tc := New(c)
			Expect(tc.Complex64()).To(Equal(c))
		})
		Specify("from ptr kind", func() {
			x := complex(float32(12), float32(13))
			v := New(&x)
			Expect(v.Complex64()).Should(Equal(x))
		})
		Specify("from other kind", func() {
			sv := New("test")
			ExpectErr(sv.Complex64()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with Complex128()", func() {
		Specify("from int*, uint*, float*, complex* kind", func() {
			c := complex(float64(12), float64(13))
			r := 12
			vs := []*Value{
				New(int(r)),
				New(int8(r)),
				New(int16(r)),
				New(int32(r)),
				New(int64(r)),
				New(uint(r)),
				New(uint8(r)),
				New(uint16(r)),
				New(uint32(r)),
				New(uint64(r)),
				New(float32(r)),
				New(float64(r)),
			}
			for _, v := range vs {
				Expect(v.Complex128()).To(Equal(complex(float64(r), 0)))
			}
			tc := New(c)
			Expect(tc.Complex128()).To(Equal(c))

			c32 := complex(float32(12), float32(13))
			tc32 := New(c32)
			Expect(tc32.Complex128()).To(Equal(c))
		})
		Specify("from ptr kind", func() {
			x := complex(float64(12), float64(13))
			v := New(&x)
			Expect(v.Complex128()).Should(Equal(x))
		})
		Specify("from other kind", func() {
			sv := New("test")
			ExpectErr(sv.Complex128()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with Get()", func() {
		Specify("from map kind", func() {
			m := map[int]int{
				1: 1,
				2: 2,
			}
			v := New(m)
			for key, val := range m {
				Expect(v.MustGet(key).Int()).Should(Equal(val))
			}
			Expect(v.MustGet(3)).Should(BeNil())
		})
		Specify("from slice kind", func() {
			s := []int{1, 2}
			v := New(s)
			for idx, elem := range s {
				Expect(v.MustGet(idx).Int()).Should(Equal(elem))
			}
			Expect(v.MustGet(3)).Should(BeNil())
		})
		Specify("from array kind", func() {
			s := [3]int{1, 2}
			v := New(s)
			for idx, elem := range s {
				Expect(v.MustGet(idx).Int()).Should(Equal(elem))
			}
			Expect(v.MustGet(4)).Should(BeNil())
		})
		Specify("from struct kind", func() {
			ss := struct {
				A, B int
			}{
				A: 1,
				B: 2,
			}
			v := New(ss)
			Expect(v.MustGet("A").Int()).Should(Equal(ss.A))
			Expect(v.MustGet("B").Int()).Should(Equal(ss.B))
			Expect(v.MustGet("C")).Should(BeNil())
		})
		Specify("from ptr kind", func() {
			s := []int{1}
			v := New(&s)
			Expect(v.MustGet(0).Int()).Should(Equal(1))
		})
		Specify("from other kind", func() {
			v := New("test")
			ExpectErr(v.Get("x")).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with Map()", func() {
		Specify("from map kind", func() {
			m := map[int]int{
				1: 1,
				2: 2,
			}
			tm := New(m).MustMap()
			Expect(len(tm)).Should(Equal(len(m)))
			for tk, tv := range tm {
				Expect(tv.Int()).Should(Equal(m[tk.MustInt()]))
			}
		})
		Specify("from slice kind", func() {
			s := []int{1, 2}
			tm := New(s).MustMap()
			Expect(len(tm)).Should(Equal(len(s)))
			for idx, elem := range tm {
				Expect(elem.Int()).Should(Equal(s[idx.MustInt()]))
			}
		})
		Specify("from array kind", func() {
			s := [3]int{1, 2}
			tm := New(s).MustMap()
			Expect(len(tm)).Should(Equal(len(s)))
			for idx, elem := range tm {
				Expect(elem.Int()).Should(Equal(s[idx.MustInt()]))
			}
		})
		Specify("from struct kind", func() {
			ss := struct {
				A, B int
			}{
				A: 1,
				B: 2,
			}
			tm := New(ss).MustMap()
			Expect(len(tm)).Should(Equal(2))
			for tk, tv := range tm {
				key, err := tk.String()
				Expect(err).Should(BeNil())
				switch key {
				case "A":
					Expect(tv.Int()).Should(Equal(ss.A))
				case "B":
					Expect(tv.Int()).Should(Equal(ss.B))
				default: // must error
					Expect(tv.Int()).Should(BeNil())
				}
			}
		})
		Specify("from ptr kind", func() {
			s := []int{1, 2}
			tm := New(&s).MustMap()
			Expect(len(tm)).Should(Equal(len(s)))
			for idx, elem := range tm {
				Expect(elem.Int()).Should(Equal(s[idx.MustInt()]))
			}
		})
		Specify("from other kind", func() {
			v := New("test")
			ExpectErr(v.Map()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with Slice()", func() {
		Specify("from slice kind", func() {
			s := []int{1, 2}
			ts := New(s).MustSlice()
			Expect(len(ts)).Should(Equal(len(s)))
			for idx, elem := range ts {
				Expect(elem.Int()).Should(Equal(s[idx]))
			}
		})
		Specify("from array kind", func() {
			s := [3]int{1, 2}
			ts := New(s).MustSlice()
			Expect(len(ts)).Should(Equal(len(s)))
			for idx, elem := range ts {
				Expect(elem.Int()).Should(Equal(s[idx]))
			}
		})
		Specify("from struct kind", func() {
			ss := struct {
				A, B int
			}{
				A: 1,
				B: 2,
			}
			ts := New(ss).MustSlice()
			Expect(len(ts)).Should(Equal(2))
			Expect(ts[0].Int()).Should(Equal(ss.A))
			Expect(ts[1].Int()).Should(Equal(ss.B))
		})
		Specify("from ptr kind", func() {
			s := []int{1, 2}
			ts := New(&s).MustSlice()
			Expect(len(ts)).Should(Equal(len(s)))
			for idx, elem := range ts {
				Expect(elem.Int()).Should(Equal(s[idx]))
			}
		})
		Specify("from other kind", func() {
			v := New("test")
			ExpectErr(v.Slice()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with AList()", func() {
		Specify("from map kind", func() {
			m := map[int]int{
				1: 1,
				2: 2,
			}
			tl := New(m).MustAList()
			Expect(len(tl)).Should(Equal(len(m)))
			for _, kv := range tl {
				Expect(kv[1].Int() /* value */).Should(Equal(m[kv[0].MustInt() /* key */]))
			}
		})
		Specify("from slice kind", func() {
			s := []int{1, 2}
			tl := New(s).MustAList()
			Expect(len(tl)).Should(Equal(len(s)))
			for _, kv := range tl {
				Expect(kv[1].Int() /* value */).Should(Equal(s[kv[0].MustInt() /* idx */]))
			}
		})
		Specify("from array kind", func() {
			s := [3]int{1, 2}
			tl := New(s).MustAList()
			Expect(len(tl)).Should(Equal(len(s)))
			for _, kv := range tl {
				Expect(kv[1].Int() /* value */).Should(Equal(s[kv[0].MustInt() /* idx */]))
			}
		})
		Specify("from struct kind", func() {
			ss := struct {
				A, B int
			}{
				A: 1,
				B: 2,
			}
			tl := New(ss).MustAList()
			Expect(len(tl)).Should(Equal(2))

			Expect(tl[0][0].String() /* field name */).Should(Equal("A"))
			Expect(tl[1][0].String() /* field name */).Should(Equal("B"))

			Expect(tl[0][1].Int() /* value */).Should(Equal(ss.A))
			Expect(tl[1][1].Int() /* value */).Should(Equal(ss.B))
		})
		Specify("from ptr kind", func() {
			s := []int{1, 2}
			tl := New(&s).MustAList()
			Expect(len(tl)).Should(Equal(len(s)))
			for _, kv := range tl {
				Expect(kv[1].Int() /* value */).Should(Equal(s[kv[0].MustInt() /* idx */]))
			}
		})
		Specify("from other kind", func() {
			v := New("test")
			ExpectErr(v.AList()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with PList()", func() {
		Specify("from map kind", func() {
			m := map[int]int{
				1: 1,
				2: 2,
			}
			tl := New(m).MustPList()
			Expect(len(tl)).Should(Equal(2 * len(m)))
			for i := 0; i < len(tl)/2; i += 2 {
				k := tl[i]
				v := tl[i+1]
				Expect(v.Int()).Should(Equal(m[k.MustInt()]))
			}
		})
		Specify("from slice kind", func() {
			s := []int{1, 2}
			tl := New(s).MustPList()
			Expect(len(tl)).Should(Equal(2 * len(s)))
			for i := 0; i < len(tl); i += 2 {
				k := tl[i]
				v := tl[i+1]
				Expect(v.Int()).Should(Equal(s[k.MustInt()]))
			}
		})
		Specify("from array kind", func() {
			s := [3]int{1, 2}
			tl := New(s).MustPList()
			Expect(len(tl)).Should(Equal(2 * len(s)))
			for i := 0; i < len(tl); i += 2 {
				k := tl[i]
				v := tl[i+1]
				Expect(v.Int()).Should(Equal(s[k.MustInt()]))
			}
		})
		Specify("from struct kind", func() {
			ss := struct {
				A, B int
			}{
				A: 1,
				B: 2,
			}
			tl := New(ss).MustPList()
			Expect(len(tl)).Should(Equal(2 * 2))
			Expect(tl[0].String() /* field name */).Should(Equal("A"))
			Expect(tl[1].Int() /* value */).Should(Equal(ss.A))
			Expect(tl[2].String() /* field name */).Should(Equal("B"))
			Expect(tl[3].Int() /* value */).Should(Equal(ss.B))
		})
		Specify("from ptr kind", func() {
			s := []int{1, 2}
			tl := New(&s).MustPList()
			Expect(len(tl)).Should(Equal(2 * len(s)))
			for i := 0; i < len(tl); i += 2 {
				k := tl[i]
				v := tl[i+1]
				Expect(v.Int()).Should(Equal(s[k.MustInt()]))
			}
		})
		Specify("from other kind", func() {
			v := New("test")
			ExpectErr(v.PList()).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
	Context("with String()", func() {
		Specify("bool kind", func() {
			b := true
			tb := New(b)
			Expect(tb.String()).To(Equal("true"))

			b = false
			tb = New(b)
			Expect(tb.String()).To(Equal("false"))
		})
		Specify("int* kind", func() {
			x := 123
			v := New(x)
			Expect(v.String()).Should(Equal("123"))
		})
		Specify("uint* kind", func() {
			x := uint(123)
			v := New(x)
			Expect(v.String()).Should(Equal("123"))
		})
		Specify("float* kind", func() {
			x := 1.2
			v := New(x)
			Expect(v.String()).Should(Equal("1.2"))

			x = 1.2e+34
			v = New(x)
			Expect(v.String()).Should(Equal("1.2e+34"))
		})
		Specify("complex* kind", func() {
			x := 1 + 2i
			v := New(x)
			Expect(v.String()).Should(Equal("(1+2i)"))

			x = -1.2e+34i + 1.2e+34
			v = New(x)
			Expect(v.String()).Should(Equal("(1.2e+34-1.2e+34i)"))
		})
		Specify("stringer type", func() {
			var x StringerTest
			v := New(&x)
			Expect(v.String()).Should(Equal(x.String()))
		})
	})
})

type StringerTest struct{}

func (s *StringerTest) String() string {
	return "test"
}

var _ = Describe("Sets", func() {
	Context("with Set()", func() {
		Specify("int kind", func() {
			m := 123
			v := New(&m)
			err := v.Set(2)
			Expect(err).Should(BeNil())
			Expect(v.Int()).Should(Equal(2))

			p := &m
			n := 456
			v = New(&p)
			err = v.Set(&n)
			Expect(err).Should(BeNil())
			Expect(v.Int()).Should(Equal(456))
		})
		Specify("value can't set", func() {
			v := New(1)
			Expect(v.Set(2)).To(BeAssignableToTypeOf((*ErrCannotSet)(nil)))
		})
		Specify("value kind unequal", func() {
			x := 1
			vx := New(&x)
			Expect(vx.Set(1.2)).To(BeAssignableToTypeOf((*ErrTypeUnequal)(nil)))
		})
	})

	Context("with Put()", func() {
		Specify("to map kind", func() {
			x := map[string]interface{}{
				"A": 1,
				"B": "b",
			}

			vx := New(x)
			Expect(vx.Put("A", 2)).Should(BeNil())
			Expect(vx.Put("C", 1.2)).Should(BeNil())

			Expect(vx.MustGet("A").Int()).Should(Equal(2))
			Expect(vx.MustGet("B").String()).Should(Equal("b"))
			Expect(vx.MustGet("C").Float64()).Should(Equal(1.2))
		})
		Specify("to slice kind", func() {
			x := []interface{}{1, "b"}

			vx := New(x)
			Expect(vx.Put(0, 2)).Should(BeNil())
			Expect(vx.Put(2, 1.2)).Should(BeNil())

			Expect(vx.MustGet(0).Int()).Should(Equal(2))
			Expect(vx.MustGet(1).String()).Should(Equal("b"))
			Expect(vx.MustGet(2).Float64()).Should(Equal(1.2))
		})
		Specify("to array kind", func() {
			x := [3]interface{}{1, "b"}

			vx := New(&x)
			Expect(vx.Put(0, 2)).Should(BeNil())
			Expect(vx.Put(2, 1.2)).Should(BeNil())

			Expect(vx.MustGet(0).Int()).Should(Equal(2))
			Expect(vx.MustGet(1).String()).Should(Equal("b"))
			Expect(vx.MustGet(2).Float64()).Should(Equal(1.2))
		})
		Specify("to struct kind", func() {
			x := struct {
				A int
				B string
				C string
			}{1, "b", ""}

			vx := New(&x)
			Expect(vx.Put("A", 2)).Should(BeNil())
			Expect(vx.Put("C", "c")).Should(BeNil())

			Expect(vx.MustGet("A").Int()).Should(Equal(2))
			Expect(vx.MustGet("B").String()).Should(Equal("b"))
			Expect(vx.MustGet("C").String()).Should(Equal("c"))
		})
		Specify("to other kind", func() {
			vx := New("a")
			Expect(vx.Put("nil", "nil")).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))

			x := 123
			vx = New(&x)
			Expect(vx.Put("nil", "nil")).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
		})
	})
})

var _ = Describe("Dos", func() {
	Context("with EachDo()", func() {
		Specify("in map", func() {
			m := map[int]int{
				1: 1,
				2: 2,
			}
			v := New(m)
			_ = v.EachDo(func(key, val *Value) error {
				Expect(m[key.MustInt()]).Should(Equal(val.MustInt()))
				return nil
			})
		})
	})
})

var _ = Describe("ConvTo", func() {
	Specify("not ptr kind", func() {
		x := true
		var y bool

		vx := New(x)
		err := vx.ConvTo(y)
		Expect(err).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
	})
	Specify("bool kind", func() {
		x := true
		var y bool

		vx := New(x)
		err := vx.ConvTo(&y)
		Expect(err).Should(BeNil())
		Expect(y).Should(Equal(x))
	})
	Specify("int kind", func() {
		x := 123
		var y int

		vx := New(x)
		err := vx.ConvTo(&y)
		Expect(err).Should(BeNil())
		Expect(y).Should(Equal(x))
	})
	Specify("uint kind", func() {
		x := uint(123)
		y := uint(0)

		vx := New(x)
		err := vx.ConvTo(&y)
		Expect(err).Should(BeNil())
		Expect(y).Should(Equal(uint(x)))
	})
	Specify("float kind", func() {
		x := 123.4
		y := 0.0

		vx := New(x)
		err := vx.ConvTo(&y)
		Expect(err).Should(BeNil())
		Expect(y).Should(Equal(x))
	})
	Specify("complex kind", func() {
		x := 1i + 2
		y := 0i + 0

		vx := New(x)
		err := vx.ConvTo(&y)
		Expect(err).Should(BeNil())
		Expect(y).Should(Equal(x))
	})
	Specify("string kind", func() {
		x := "abc"
		y := ""

		vx := New(x)
		err := vx.ConvTo(&y)
		Expect(err).Should(BeNil())
		Expect(y).Should(Equal(x))
	})
	Specify("slice kind", func() {
		x := []interface{}{
			1, "a", 0.1,
		}

		y := make([]interface{}, 0)

		vx := New(x)
		err := vx.ConvTo(&y)
		Expect(err).Should(BeNil())
		Expect(len(y)).Should(Equal(len(x)))
		for i, e := range y {
			Expect(e).Should(Equal(x[i]))
		}
	})
	Specify("array kind", func() {
		x := [3]int{1, 2, 3}
		var y [3]int

		vx := New(x)
		err := vx.ConvTo(&y)
		Expect(err).Should(BeNil())
		for i, v := range x {
			Expect(y[i]).Should(Equal(v))
		}
	})
	Specify("map kind", func() {
		x := map[string]interface{}{
			"A": 1,
			"B": "a",
			"C": []int{1, 2, 3},
		}

		y := make(map[string]interface{})

		vx := New(x)
		err := vx.ConvTo(&y)
		Expect(err).Should(BeNil())
		for k, v := range x {
			Expect(y[k]).Should(Equal(v))
		}
	})
	Specify("to struct kind", func() {
		x := map[string]interface{}{
			"a": 1,
			"A": 11,

			"B": "b",
			"b": "bb",

			"C": "c",
			"c": "cc",
		}

		var y struct {
			A int    `value:"a"`
			B string `value:"-"`
			C string
		}

		vx := New(x)
		err := vx.ConvTo(&y)
		Expect(err).Should(BeNil())
		Expect(y.A).Should(Equal(1))
		Expect(y.B).Should(Equal(""))
		Expect(y.C).Should(Equal("c"))
	})

	Specify("time.Duration type", func() {
		x := "1s"
		var y time.Duration

		vx := New(x)
		err := vx.ConvTo(&y)
		Expect(err).Should(BeNil())
		Expect(y).Should(Equal(1 * time.Second))
	})

	Specify("time.Time type", func() {
		now := time.Now()
		x := now.Format(TimeLayout)
		var y time.Time

		vx := New(x)
		err := vx.ConvTo(&y)
		Expect(err).Should(BeNil())
		Expect(y.Format(TimeLayout)).Should(Equal(x))
	})

	Specify("net.IP type", func() {
		x := "8.8.8.8"
		var y net.IP

		vx := New(x)
		err := vx.ConvTo(&y)
		Expect(err).Should(BeNil())
		Expect(y).Should(Equal(net.ParseIP(x)))
	})

	Specify("url.URL type", func() {
		x := "http://userinfo@www.abc.com/path?query=1#fragment"
		var y url.URL

		vx := New(x)
		err := vx.ConvTo(&y)
		Expect(err).Should(BeNil())
		u, _ := url.Parse(x)
		Expect(y).Should(Equal(*u))
	})

	Specify("mail.Address", func() {
		x := "name <name@host.com>"
		var y mail.Address

		vx := New(x)
		err := vx.ConvTo(&y)
		Expect(err).Should(BeNil())
		a, _ := mail.ParseAddress(x)
		Expect(y).Should(Equal(*a))
	})

	Specify("regexp.Regexp", func() {
		x := "[0-9]+"
		var y regexp.Regexp

		vx := New(x)
		err := vx.ConvTo(&y)
		Expect(err).Should(BeNil())
		r, _ := regexp.Compile(x)
		Expect(y).Should(Equal(*r))
	})

	Specify("ByteSize type", func() {
		x := "10KB"
		var y ByteSize

		vx := New(x)
		err := vx.ConvTo(&y)
		Expect(err).Should(BeNil())
		Expect(y).Should(Equal(ByteSize(10 * 1024)))
	})

	Specify("nest struct kind", func() {
		type xx struct {
			X int
			Y int
		}
		var y struct {
			A int
			B time.Duration
			C *xx
		}
		x := map[string]interface{}{
			"A": 1,
			"B": "1s",
			"C": map[string]int{
				"X": 10,
				"Y": 11,
			},
		}

		vx := New(x)
		err := vx.ConvTo(&y)
		Expect(err).Should(BeNil())
		Expect(y.A).Should(Equal(1))
		Expect(y.B).Should(Equal(time.Second))
		Expect(y.C.X).Should(Equal(10))
		Expect(y.C.Y).Should(Equal(11))
	})
	Specify("to chan kind", func() {
		x := map[string]interface{}{
			"A": 1,
			"B": "a",
		}
		var y chan int

		vx := New(x)
		err := vx.ConvTo(&y)
		Expect(err).To(BeAssignableToTypeOf((*ErrUnsupportedKind)(nil)))
	})
})

var _ = Describe("Musts", func() {
	Specify("with MustInt8()", func() {
		x := int8(1)
		v := New(x)
		Expect(v.MustInt8()).To(Equal(x))

		v = New("test")
		Expect(func() { v.MustInt8() }).Should(Panic())
	})
	Specify("with MustInt16()", func() {
		x := int16(1)
		v := New(x)
		Expect(v.MustInt16()).To(Equal(x))

		v = New("test")
		Expect(func() { v.MustInt16() }).Should(Panic())
	})
	Specify("with MustInt32()", func() {
		x := int32(1)
		v := New(x)
		Expect(v.MustInt32()).To(Equal(x))

		v = New("test")
		Expect(func() { v.MustInt32() }).Should(Panic())
	})
	Specify("with MustInt64()", func() {
		x := int64(1)
		v := New(x)
		Expect(v.MustInt64()).To(Equal(x))

		v = New("test")
		Expect(func() { v.MustInt64() }).Should(Panic())
	})
	Specify("with MustInt()", func() {
		x := int(1)
		v := New(x)
		Expect(v.MustInt()).To(Equal(x))

		v = New("test")
		Expect(func() { v.MustInt() }).Should(Panic())
	})
	Specify("with MustFloat32()", func() {
		x := float32(1)
		v := New(x)
		Expect(v.MustFloat32()).To(Equal(x))

		v = New("test")
		Expect(func() { v.MustFloat32() }).Should(Panic())
	})
	Specify("with MustFloat64()", func() {
		x := float64(1)
		v := New(x)
		Expect(v.MustFloat64()).To(Equal(x))

		v = New("test")
		Expect(func() { v.MustFloat64() }).Should(Panic())
	})
	Specify("with MustGet()", func() {
		x := map[int]int{1: 1}
		v := New(x)
		Expect(v.MustGet(1).Int()).To(Equal(1))

		v = New("test")
		Expect(func() { v.MustGet(1) }).Should(Panic())
	})
	Specify("with MustMap()", func() {
		x := map[int]int{1: 1}
		v := New(x)
		m := v.MustMap()
		for k, v := range m {
			Expect(v.MustInt()).To(Equal(x[k.MustInt()]))
		}

		v = New("test")
		Expect(func() { v.MustMap() }).Should(Panic())
	})
	Specify("with MustSlice()", func() {
		x := []int{1, 2}
		v := New(x)
		s := v.MustSlice()
		for i, v := range s {
			Expect(v.MustInt()).To(Equal(x[i]))
		}

		v = New("test")
		Expect(func() { v.MustSlice() }).Should(Panic())
	})
	Specify("with MustAList()", func() {
		x := map[int]int{1: 1}
		v := New(x)
		m := v.MustAList()
		for _, kv := range m {
			Expect(kv[1].MustInt()).To(Equal(x[kv[0].MustInt()]))
		}

		v = New("test")
		Expect(func() { v.MustAList() }).Should(Panic())
	})
	Specify("with MustPList()", func() {
		x := map[int]int{1: 1}
		v := New(x)
		p := v.MustPList()
		Expect(p[1].MustInt()).To(Equal(x[p[0].MustInt()]))

		v = New("test")
		Expect(func() { v.MustPList() }).Should(Panic())
	})
})

var _ = Describe("Errs", func() {
	Specify("of ErrNumOverflow", func() {
		m := "method"
		k := reflect.Int
		es := "table: call of " + m + " overflows " + k.String()
		Expect((&ErrNumOverflow{m, k}).Error()).To(Equal(es))
	})
	Specify("of ErrCannotBeNil", func() {
		m := "method"
		es := "table: call of " + m + " on nil value"
		Expect((&ErrCannotBeNil{m}).Error()).To(Equal(es))
	})
	Specify("of ErrNotExist", func() {
		m := "method"
		k := "Int"
		es := "table: call of " + m + " not exist of " + k
		Expect((&ErrNotExist{m, k}).Error()).To(Equal(es))
	})
	Specify("of ErrCannotSet", func() {
		m := "method"
		es := "table: call of " + m + " on unaddressable value"
		Expect((&ErrCannotSet{m}).Error()).To(Equal(es))
	})
	Specify("of ErrNumOverflow", func() {
		m := "method"
		k := reflect.Int
		es := "table: call of " + m + " overflows " + k.String()
		Expect((&ErrNumOverflow{m, k}).Error()).To(Equal(es))
	})
	Specify("of ErrTypeUnequal", func() {
		m := "method"
		k1 := reflect.Int
		k2 := reflect.Float32
		es := "table: call of " + m + " between " + k1.String() + " and " + k2.String()
		Expect((&ErrTypeUnequal{m, k1, k2}).Error()).To(Equal(es))
	})
	Specify("of ErrOutOfRange", func() {
		m := "method"
		es := "table: call of " + m + " out of range"
		Expect((&ErrOutOfRange{m}).Error()).To(Equal(es))
	})
	Specify("of ErrUnsupportedKind", func() {
		m := "method"
		k := reflect.Int
		es := "table: call of " + m + " on " + k.String() + " value"
		Expect((&ErrUnsupportedKind{m, k}).Error()).To(Equal(es))

		ks := "type"
		es = "table: call of " + m + " on " + ks + " value"
		Expect((&ErrUnsupportedKind{m, ks}).Error()).To(Equal(es))
	})
})
