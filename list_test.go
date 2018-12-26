package fmr

import (
	"testing"
)

func TestList(t *testing.T) {
	//Debug = true
	cases := []string{
		`北京`,
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

		t.Logf("%+v\n", p)
		states := p.GetFinalStates()
		if len(states) > 0 {
			t.Logf("%s\n%+v\n", c, states)
			for _, f := range states {
				t.Log(f)
				trees := p.GetTrees(f)
				t.Log(trees)
				for _, tree := range trees {
					sem, err := tree.Semantic()
					if err != nil {
						t.Error(err)
					}
					t.Log(sem)
				}
			}
		} else {
			t.Logf("%s\nno result\n", c)
		}
	}
}
