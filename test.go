package main

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

type TestParam struct {
	C string `json:"c,omitempty"`
	A string `json:"a,omitempty"`
	B string `json:"b,omitempty"`
}

func main() {
	p := TestParam{
		A: "a",
		C: "c",
	}

	b, _ := jsoniter.Marshal(&p)
	fmt.Println(string(b))

	p2 := TestParam{}
	_ = jsoniter.Unmarshal(b, &p2)
	fmt.Printf("%+v\n", p2)
}
