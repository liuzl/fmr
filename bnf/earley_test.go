package bnf

import (
	"os"
	"strings"
	"testing"
)

func NewParser(g *Grammar, start, text string) *Parser {
	tokens := strings.Fields(text)
	parser := &Parser{g: g}
	parser.columns = append(parser.columns, &TableColumn{index: 0, token: ""})
	for i, token := range tokens {
		parser.columns = append(parser.columns,
			&TableColumn{index: i + 1, token: token})
	}
	parser.finalState = parser.parse(start)
	return parser
}

func TestEarleyParse(t *testing.T) {
	grammar := `<sym> = "a";
	<op> = "+";
	<expr> = <sym> | <expr> <op> <expr>;`
	grammar = `<expr> = "a" | "a" "+" <expr>;`
	//grammar = `<expr> = "a";`
	strs := []string{
		"a",
		"a + a",
		"a + a + a",
		"a + a + a + a",
		//"a + a + a + a + a",
		//"a + a + a + a + a + a",
		//"a + a + a + a + a + a + a",
		"+ a",
	}
	g, err := Parse(grammar)
	if err != nil {
		t.Error(err)
	}
	_, err = JsonMarshalIndent(g, "", " ")
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
		}
	}
}
