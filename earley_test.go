package fmr

import (
	"bytes"
	"testing"

	"zliu.org/goutil"
)

func TestEarleyParse(t *testing.T) {
	var grammar = `<expr> = "a" | "a" "+" <expr> {nf.math.sum($1, $3)};`
	//grammar = `<expr> = "a";`
	strs := []string{
		"a",
		"a + a",
		//"a + a + a",
		//"a + a + a + a",
		//"a + a + a + a + a",
		//"a + a + a + a + a + a",
		//"a + a + a + a + a + a + a",
		"+ a",
	}
	g, err := GrammarFromString(grammar, "a")
	if err != nil {
		t.Error(err)
	}
	_, err = goutil.JSONMarshalIndent(g, "", " ")
	if err != nil {
		t.Error(err)
	}
	//fmt.Println(string(b))
	for _, text := range strs {
		p, err := g.EarleyParse(text, "expr")
		if err != nil {
			t.Error(err)
		}
		t.Logf("%+v\n", p)
		for _, finalState := range p.finalStates {
			trees := p.GetTrees(finalState)
			t.Log("tree number:", len(trees))
			for _, tree := range trees {
				var buf bytes.Buffer
				tree.Print(&buf)
				t.Log(buf.String())
				tree.TreePrint()
				b, err := goutil.JSONMarshalIndent(tree, "", " ")
				if err != nil {
					t.Error(err)
				}
				t.Logf("%+v", string(b))
			}
		}
	}
}
