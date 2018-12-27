package fmr

import (
	"bytes"
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
					var buf bytes.Buffer
					tree.Print(&buf)
					t.Log(buf.String())
					sem, err := tree.Semantic()
					if err != nil {
						t.Error(err)
					}
					t.Log(sem)
					s, err := tree.Eval()
					if err != nil {
						t.Error(err)
					}
					t.Log(s)
				}
			}
		} else {
			t.Logf("%s\nno result\n", c)
		}
	}
}
