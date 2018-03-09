package bnf

import (
	//"fmt"
	"testing"
)

var tests = []string{
	`<list>  =  "<" <items> ">"               ;
	<items> =  <items> " " <item> {     nf.math.sum($1,$3)} | <item>   ;
	<item>  =  "f    \\uoo\n" | "bar\t" | "baz"|"好吧"         ;
	`,
}

func TestLex(t *testing.T) {
	for _, c := range tests {
		g, err := Parse(c)
		if err != nil {
			t.Error(err)
		}
		_, err = JsonMarshalIndent(g, "", "  ")
		if err != nil {
			t.Error(err)
		}
		//fmt.Println(string(b))
	}
}
