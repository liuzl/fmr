package earley

import (
	"os"
	"testing"
)

func TestEarleyParse(t *testing.T) {
	SYM := NewRule("SYM", NewProduction(&Terminal{"a"}))
	OP := NewRule("OP", NewProduction("+"))
	EXPR := NewRule("EXPR", NewProduction(SYM))
	EXPR.Add(NewProduction(EXPR, OP, EXPR))

	strs := []string{
		//"a",
		"a + a",
		"a + a + a",
		//"a + a + a + a",
		//"a + a + a + a + a",
		//"a + a + a + a + a + a",
		//"a + a + a + a + a + a + a",
		"+ a",
	}
	for _, text := range strs {
		p := NewParser(EXPR, text)
		trees := p.GetTrees()
		t.Log("tree number:", len(trees))
		for _, tree := range trees {
			tree.Print(os.Stdout)
		}
	}
}
