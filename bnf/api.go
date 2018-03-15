package bnf

import (
	"fmt"
	"github.com/liuzl/ling"
	"strings"
)

var nlp = ling.MustNLP(ling.Norm)

// EarleyParse parses text for rule <start>
func (g *Grammar) EarleyParse(start, text string) (*Parse, error) {
	if start = strings.TrimSpace(start); start == "" {
		return nil, fmt.Errorf("start rule is empty")
	}
	if g.Rules[start] == nil {
		return nil, fmt.Errorf("start rule:<%s> not found in Grammar", start)
	}
	if text = strings.TrimSpace(text); text == "" {
		return nil, fmt.Errorf("text is empty")
	}
	d := ling.NewDocument(text)
	if err := nlp.Annotate(d); err != nil {
		return nil, err
	}

	parse := &Parse{g: g}
	parse.columns = append(parse.columns, &TableColumn{index: 0, token: ""})
	for _, token := range d.Tokens {
		if token.Type == ling.Space {
			continue
		}
		parse.columns = append(parse.columns,
			&TableColumn{
				index: len(parse.columns),
				token: token.Annotations[ling.Norm],
			})
	}

	parse.finalState = parse.parse(start)
	return parse, nil
}
