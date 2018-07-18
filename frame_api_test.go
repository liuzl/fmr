package fmr

import (
	"io/ioutil"
	"testing"
)

func TestMatchFrames(t *testing.T) {
	cases := []string{
		`从北京飞上海`,
	}
	b, err := ioutil.ReadFile("sf.grammar")
	if err != nil {
		t.Error(err)
	}
	g, err := CFGrammar(string(b))
	if err != nil {
		t.Error(err)
	}
	for _, c := range cases {
		if err = g.MatchFrames(c); err != nil {
			t.Error(err)
		}
	}
}
