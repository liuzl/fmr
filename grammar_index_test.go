package fmr

import (
	"io/ioutil"
	"testing"
)

func TestGrammarIndex(t *testing.T) {
	b, err := ioutil.ReadFile("sf.grammar")
	if err != nil {
		t.Error(err)
	}
	g, err := CFGrammar(string(b))
	if err != nil {
		t.Error(err)
	}
	t.Log(g)
}
