package fmr

import (
	"bytes"
	"testing"
)

func TestList(t *testing.T) {
	//Debug = true
	cases := []string{
		`直辖市：北京上海天津`,
		`直辖市：北京、上海和天津`,
		`直辖市：北京、上海和天津、津城`,
		`直辖市：帝都、魔都、寨都、旧都`,
		`直辖市：北京`,
	}
	g, err := GrammarFromFile("sf.grammar")
	if err != nil {
		t.Error(err)
	}
	for _, c := range cases {
		t.Log(c)
		trees, err := g.Parse(c, "cities")
		if err != nil {
			t.Error(err)
		}
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
}
