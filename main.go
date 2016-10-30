package main

import (
	"log"
	"zjson/core"
)

func main() {
	log.Println("hello zjson")

	t1 := `"abc"`
	t2 := `{"a":"12"}`
	t3 := `{"abc":["10"]}`
	t4 := `{"abc": "tt","hh": "h1","cc": {}}`
	t5 := `{"a":12,"b":25}`
	t6 := `{}`
	t7 := `{"abs":{"yy":[{"aa":"123"}, "123", 123]}}`
	t8 := `{
		"abs":{
			"yy":[
				{"aa":"123"},
				"123",
				123
			]
		}
	}`
	_ = t1
	_ = t2
	_ = t3
	_ = t4
	_ = t5
	_ = t6
	_ = t7
	_ = t8

	ts := core.NewJSONTokenSpliter()

	tokens, err := ts.Run([]byte(t8))
	if err != nil {
		panic(err)
	}
	log.Println("tokens", tokens)

	p := core.NewJSONParser()
	j := p.Parse(tokens)
	log.Println("json", j)
}
