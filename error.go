package value

import (
	"reflect"
)

type (
	// ErrNumOverflow ...
	ErrNumOverflow struct {
		Method string
		Kind   reflect.Kind
	}

	// ErrUnsupportedKind ...
	ErrUnsupportedKind struct {
		Method string
		Kind   interface{}
	}

	// ErrCannotBeNil ...
	ErrCannotBeNil struct {
		Method string
	}

	// ErrNotExist ...
	ErrNotExist struct {
		Method string
		Thing  string
	}

	// ErrCannotSet ...
	ErrCannotSet struct {
		Method string
	}

	// ErrTypeUnequal ...
	ErrTypeUnequal struct {
		Method string
		Kind1  reflect.Kind
		Kind2  reflect.Kind
	}

	// ErrOutOfRange ...
	ErrOutOfRange struct {
		Method string
	}
)

func (e *ErrUnsupportedKind) Error() string {
	rkind, ok := e.Kind.(reflect.Kind)
	if ok && rkind == 0 {
		return "table: call of " + e.Method + " on zero value"
	}

	var kind string
	if ok {
		kind = rkind.String()
	} else {
		kind, _ = e.Kind.(string)
	}

	return "table: call of " + e.Method + " on " + kind + " value"
}

func (e *ErrNumOverflow) Error() string {
	return "table: call of " + e.Method + " overflows " + e.Kind.String()
}

func (e *ErrCannotBeNil) Error() string {
	return "table: call of " + e.Method + " on nil value"
}

func (e *ErrNotExist) Error() string {
	return "table: call of " + e.Method + " not exist of " + e.Thing
}

func (e *ErrCannotSet) Error() string {
	return "table: call of " + e.Method + " on unaddressable value"
}

func (e *ErrTypeUnequal) Error() string {
	return "table: call of " + e.Method + " between " + e.Kind1.String() + " and " + e.Kind2.String()
}

func (e *ErrOutOfRange) Error() string {
	return "table: call of " + e.Method + " out of range"
}
