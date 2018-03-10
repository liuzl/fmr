package bnf

import (
	"strings"
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
