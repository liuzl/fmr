package fmr

import (
	"io/ioutil"
	"testing"
)

func TestMatchFrames(t *testing.T) {
	cases := []string{
		`从北京飞上海`,
		`飞上海 从北京`,
		`我要从北京走`,
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
		if frames, err := g.MatchFrames(c); err != nil {
			t.Error(err)
		} else {
			for k, sf := range frames {
				t.Log(k)
				t.Log(sf)
				for term, slots := range sf.Fillings {
					t.Log(term)
					for _, slot := range slots {
						for _, tree := range slot.Trees {
							t.Log(tree.Semantic())
						}
					}
				}
			}
		}
		if err := g.FrameFMR(c); err != nil {
			t.Error(err)
		}
	}
}
