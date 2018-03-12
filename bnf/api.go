package bnf

import (
	"github.com/liuzl/ling"
)

var nlp = ling.MustNLP(ling.Norm)

func NewParser(g *Grammar, start, text string) *Parser {
	d := ling.NewDocument(text)
	err := nlp.Annotate(d)
	if err != nil {
		//TODO
		return nil
	}
	var tokens []string
	for _, token := range d.Tokens {
		if token.Type == ling.Space {
			continue
		}
		tokens = append(tokens, token.Annotations[ling.Norm])
	}
	parser := &Parser{g: g}
	parser.columns = append(parser.columns, &TableColumn{index: 0, token: ""})
	for i, token := range tokens {
		parser.columns = append(parser.columns,
			&TableColumn{index: i + 1, token: token})
	}
	parser.finalState = parser.parse(start)
	return parser
}
