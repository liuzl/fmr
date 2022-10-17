package fmr

import (
	//"fmt"
	"testing"

	"zliu.org/goutil"
)

var tests = []string{
	`<list>  =  "<" <items> ">"               ;
	<items> =  <items> " " <item> {     nf.math.sum($1,$3)} | <item>   ;
	<item>  =  "(?ilfw)f    \\uoo\n" | "bar\t" | "baz"|"好吧"         ;
	<name> = "\(" (any) ")" ;
	`,
	`<datetimes> = (list<datetime>);
	<datetime>="20181219"|"20181218";
	`,
}

func TestLex(t *testing.T) {
	for _, c := range tests {
		g, err := GrammarFromString(c, "test")
		if err != nil {
			t.Error(err)
		}
		b, err := goutil.JSONMarshalIndent(g, "", "  ")
		if err != nil {
			t.Error(err)
		}
		t.Log(string(b))
	}
}
