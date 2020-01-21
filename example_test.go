package value_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"time"
	"net/mail"

	"github.com/helloyi/go-value"
)

func Example() {
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

func ExampleValue_EachDo_map() {
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

func ExampleValue_ConvTo() {
	x := map[string]interface{}{
		"a": 1,
		"B": "b",
		"C": "c",
		"D": "13s",
		"E": "Fri Nov 1 19:13:55 +0800 CST 2019",
		"F": "8.8.8.8",
		"G": "http://host/path?param=x",
		"H": "[0-9]+",
		"I": "name <uer@mail.com>",
	}

	var y struct {
		A int    `value:"a"` // set with name "a"
		B int    `value:"_"` // passed
		C string // set with name "C"
		D time.Duration
		E *time.Time
		F *net.IP
		G *url.URL
		H *regexp.Regexp
		I *mail.Address
	}

	v := value.New(x)
	if err := v.ConvTo(&y); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%v\n", y)
	// Output: {1 0 c 13s 2019-11-01 19:13:55 +0800 CST 8.8.8.8 http://host/path?param=x [0-9]+ "name" <uer@mail.com>}
}

func ExampleValue_Put() {
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
