package bnf

import (
	"fmt"
	"github.com/liuzl/ling"
	"strings"
)

var nlp = ling.MustNLP(ling.Norm)

// EarleyParse parses text for rule <start>
func (g *Grammar) EarleyParse(start, text string) (*Parse, error) {
	tokens, err := getTokens(text)
	if err != nil {
		return nil, err
	}
	return g.earleyParse(start, tokens)
}

// EarleyParseAll extracts all submatches in text for rule <start>
func (g *Grammar) EarleyParseAll(start, text string) ([]*Parse, error) {
	tokens, err := getTokens(text)
	if err != nil {
		return nil, err
	}
	var ret []*Parse
	for i := 0; i < len(tokens)-1; {
		p, err := g.earleyParse(start, tokens[i:])
		if err != nil {
			return nil, err
		}
		if p.finalState != nil {
			ret = append(ret, p)
			i += p.finalState.End
		} else {
			i++
		}
	}
	return ret, nil
}

func (g *Grammar) earleyParse(start string, tokens []*ling.Token) (*Parse, error) {
	if start = strings.TrimSpace(start); start == "" {
		return nil, fmt.Errorf("start rule is empty")
	}
	if g.Rules[start] == nil {
		return nil, fmt.Errorf("start rule:<%s> not found in Grammar", start)
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens to parse")
	}

	parse := &Parse{g: g}
	parse.columns = append(parse.columns, &TableColumn{index: 0, token: ""})
	for _, token := range tokens {
		parse.columns = append(parse.columns,
			&TableColumn{
				index: len(parse.columns),
				token: token.Annotations[ling.Norm],
			})
	}

	parse.finalState = parse.parse(start)
	return parse, nil
}

func getTokens(text string) ([]*ling.Token, error) {
	if text = strings.TrimSpace(text); text == "" {
		return nil, fmt.Errorf("text is empty")
	}
	d := ling.NewDocument(text)
	if err := nlp.Annotate(d); err != nil {
		return nil, err
	}
	var ret []*ling.Token
	for _, token := range d.Tokens {
		if token.Type == ling.Space {
			continue
		}
		ret = append(ret, token)
	}
	return ret, nil
}
