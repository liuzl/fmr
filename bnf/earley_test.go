package bnf

import (
	"github.com/liuzl/goutil"
	"os"
	"testing"
)

func TestEarleyParse(t *testing.T) {
	grammar := `<sym> = "a";
	<op> = "+";
	<expr> = <sym> | <expr> <op> <expr>;`
	grammar = `<expr> = "a" | "a" "+" <expr> {nf.math.sum($1, $3)};`
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
	g, err := CFGrammar(grammar)
	if err != nil {
		t.Error(err)
	}
	_, err = goutil.JsonMarshalIndent(g, "", " ")
	if err != nil {
		t.Error(err)
	}
	//fmt.Println(string(b))
	for _, text := range strs {
		p := NewParser(g, "expr", text)
		t.Logf("%+v\n", p)
		trees := p.GetTrees()
		t.Log("tree number:", len(trees))
		for _, tree := range trees {
			tree.Print(os.Stdout)
			b, err := goutil.JsonMarshalIndent(tree, "", " ")
			if err != nil {
				t.Error(err)
			}
			t.Logf("%+v", string(b))
		}
	}
}
