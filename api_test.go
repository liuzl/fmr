package fmr

import (
	"testing"
)

func TestLocalParse(t *testing.T) {
	tests := []string{
		`柏乡位于河北省`,
	}
	g := &Grammar{}
	for _, c := range tests {
		ps, err := g.EarleyParseMaxAll(c, "loc_province", "loc_county")
		if err != nil {
			t.Error(err)
		}
		for _, p := range ps {
			for _, f := range p.GetFinalStates() {
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
		}
	}
}

func TestGParse(t *testing.T) {
	tests := []string{
		`柏乡位于河北省`,
	}
	grammar := `<loc> = <loc_province> {nf.loc($1)}| <loc_county> {nf.loc($1)};`
	g, err := GrammarFromString(grammar, "loc")
	if err != nil {
		t.Error(err)
	}
	for _, c := range tests {
		ps, err := g.EarleyParseMaxAll(c, "loc")
		if err != nil {
			t.Error(err)
		}
		for _, p := range ps {
			for _, f := range p.GetFinalStates() {
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
		}
	}
}
