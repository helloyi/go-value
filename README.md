# The Value in go

[![GoDoc](https://godoc.org/github.com/helloyi/go-value?status.svg)](https://godoc.org/github.com/helloyi/go-value) [![Build Status](https://travis-ci.com/helloyi/go-value.svg?branch=master)](https://travis-ci.com/helloyi/go-value) [![Coverage Status](https://coveralls.io/repos/github/helloyi/go-value/badge.svg?branch=master)](https://coveralls.io/github/helloyi/go-value?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/helloyi/go-value)](https://goreportcard.com/report/github.com/helloyi/go-value) ![stability-unstable](https://img.shields.io/badge/stability-unstable-yellow.svg) [![License](https://img.shields.io/github/license/helloyi/go-value)](https://github.com/helloyi/go-value/blob/master/LICENSE)

Manipulates(gets,sets,convert etc.) interface{} value, convenient and simple, and reflect-based.

## Features

You can do the following things, after new Value from interface{} like this `v := value.New(var)`.
+ Get value from map slice array or struct with key idx or fieldname
+ Put key value pair to map slice or array
+ Set value
+ For range value any type
+ Get value with Int*, Float*, PList etc.
+ Numerical type is adaptive when get; example, `i` kind is uint8 and get with Int16() is legal, the value is automatically converted to int16 instead of type mismatch.
+ Provide `Must*` API, for chaining call and some friendly writing.
+ Unmarshal to a value with `Value.ConvTo`
+ Supported covert to time.Duration or time.Time etc..
+ [ ] Supported convert to frequently used build-in type

## Example

[example_test.go](https://github.com/helloyi/go-value/blob/master/example_test.go)

convert v to net.IP, time.Duration etc..
```golang
package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/helloyi/go-value"
)

func main() {
	x := map[string]interface{}{
		"a": 1,
		"B": "b",
		"C": "c",
		"D": "13s",
		"E": "Fri Nov 1 19:13:55 +0800 CST 2019",
		"F": "8.8.8.8",
	}

	var y struct {
		A int    `value:"a"` // set with name "a"
		B int    `value:"_"` // passed
		C string // set with name "C"
		D time.Duration
		E *time.Time
		F *net.IP
	}

	v := value.New(x)
	if err := v.ConvTo(&y); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%v\n", y)
	// Output: {1 0 c 13s 2019-11-01 19:13:55 +0800 CST 8.8.8.8}
}
```
