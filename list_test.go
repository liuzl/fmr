package fmr

import (
	"testing"
)

func TestList(t *testing.T) {
	Debug = true
	cases := []string{
		`北京上海天津`,
	}
	g, err := GrammarFromFile("sf.grammar")
	if err != nil {
		t.Error(err)
	}
	for _, c := range cases {
		p, err := g.EarleyParse(c, "cities")
		if err != nil {
			t.Error(err)
		}
		states := p.GetFinalStates()
		if len(states) > 0 {
			t.Logf("%s\n%+v\n", c, states)
		} else {
			t.Logf("%s\nno result\n", c)
		}
	}
}
