package fmr

import (
	"testing"
)

func TestGrammarIndex(t *testing.T) {
	g, err := GrammarFromFile("sf.grammar")
	if err != nil {
		t.Error(err)
	}
	t.Log(g)
}
