package fmr

import (
	"testing"
)

func TestMatchFrames(t *testing.T) {
	cases := []string{
		`从北京飞上海`,
		`飞上海，从北京，后天`,
		`我要从北京走`,
	}
	g, err := GrammarFromFile("sf.grammar")
	if err != nil {
		t.Error(err)
	}
	for _, c := range cases {
		fmrs, err := g.FrameFMR(c)
		if err != nil {
			t.Error(err)
		}
		t.Log(c, fmrs)
	}
}
