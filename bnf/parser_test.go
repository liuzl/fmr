package bnf

import (
	"testing"
)

func TestParse(t *testing.T) {
	for _, c := range tests {
		g, err := Parse(c)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%+v", g)
	}
}
