package value_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/helloyi/go-value"
)

func ExampleGets() {
	resp, err := http.Get("https://api.alternative.me/fng/?limit=10")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var res interface{}
	if err := json.Unmarshal(data, &res); err != nil {
		log.Fatalln(err)
	}

	rest := value.New(res)

	fngs := map[string]string{}
	_ = rest.MustGet("data").EachDo(func(_, fng *value.Value) error {
		ts, _ := fng.MustGet("timestamp").String()
		val, _ := fng.MustGet("value").String()
		fngs[ts] = val
		return nil
	})

	data, err = json.MarshalIndent(fngs, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(data))
}

func ExampleForMap() {
	var m interface{} = map[int]int{
		1: 1,
		2: 2,
	}
	mv := value.New(m)

	sum := 0
	if err := mv.EachDo(func(key, val *value.Value) error {
		sum += val.MustInt()
		return nil
	}); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(sum)
	// Output: 3
}

func ExampleConvTo() {
	x := map[string]interface{}{
		"a": 1,
		"A": 11,

		"B": "b",
		"b": "bb",

		"C": "c",
		"c": "cc",
	}

	var y struct {
		A int    `value:"a"` // set with name "a"
		B int    `value:"_"` // passed
		C string // set with name "C"
	}

	v := value.New(x)
	if err := v.ConvTo(&y); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%v\n", y)
	// Output: {1 0 c}
}

func ExampleSetMap() {
	var m interface{} = map[int]int{
		1: 1,
		2: 2,
	}
	mv := value.New(m)

	if err := mv.Put(1, 100); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(mv.MustGet(1).MustInt())
	// Output: 100
}
