package bnf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

var tests = []string{
	`<list>  =  "<" <items> ">"               ;
	<items> =  <items> " " <item> {     nf.math.sum($1,$3)} | <item>   ;
	<item>  =  "f    \\uoo\n" | "bar\t" | "baz"|"好吧"         ;
	`,
}

func JsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

func JsonMarshalIndent(t interface{}, prefix, indent string) ([]byte, error) {
	b, err := JsonMarshal(t)
	if err != nil {
		return b, err
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, prefix, indent)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func TestLex(t *testing.T) {
	for _, c := range tests {
		g, err := Parse(c)
		if err != nil {
			t.Error(err)
		}
		b, err := JsonMarshalIndent(g, "", "  ")
		if err != nil {
			t.Error(err)
		}
		fmt.Println(string(b))
	}
}
