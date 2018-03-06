package bnf

import (
	"testing"
)

var tests = []string{
	`query_fee {
		want = "想" | "要"
		query = "查" | "查\\\\询"
		fee = "\"手机费" | "话费"
		task = ("我")? (want)? query fee 
	}`,
}

func TestLex(t *testing.T) {
	for _, c := range tests {
		l := lex(c)
		for token := range l.items {
			t.Logf("%+v", token.String())
		}
	}
}
