package fmr

import (
	"os"
	"testing"
)

func TestList(t *testing.T) {
	//Debug = true
	cases := []string{
		`直辖市：北京上海天津`,
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
			for _, f := range states {
				trees := p.GetTrees(f)
				for _, tree := range trees {
					tree.Print(os.Stdout)
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
