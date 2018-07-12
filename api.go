package fmr

import (
	"fmt"
	"strings"
	"sync"

	"github.com/liuzl/ling"
)

var nlp *ling.Pipeline
var once sync.Once

func NLP() *ling.Pipeline {
	once.Do(func() {
		var err error
		var tagger *ling.DictTagger
		if nlp, err = ling.DefaultNLP(); err != nil {
			panic(err)
		}
		if tagger, err = ling.NewDictTagger(); err != nil {
			panic(err)
		}
		if err = nlp.AddTagger(tagger); err != nil {
			panic(err)
		}
	})
	return nlp
}

// EarleyParse parses text for rule <start>
func (g *Grammar) EarleyParse(start, text string) (*Parse, error) {
	tokens, l, err := extract(text)
	if err != nil {
		return nil, err
	}
	return g.earleyParse(start, text, tokens, l)
}

// EarleyParseAll extracts all submatches in text for rule <start>
func (g *Grammar) EarleyParseAll(start, text string) ([]*Parse, error) {
	tokens, l, err := extract(text)
	if err != nil {
		return nil, err
	}
	var ret []*Parse
	for i := 0; i < len(tokens); {
		p, err := g.earleyParse(start, text, tokens[i:], l)
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

func (g *Grammar) earleyParse(start, text string,
	tokens []*ling.Token, l *Grammar) (*Parse, error) {
	if start = strings.TrimSpace(start); start == "" {
		return nil, fmt.Errorf("start rule is empty")
	}
	if g.Rules[start] == nil {
		return nil, fmt.Errorf("start rule:<%s> not found in Grammar", start)
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens to parse")
	}

	parse := &Parse{grammars: []*Grammar{g}, text: text}
	if l != nil {
		parse.grammars = append(parse.grammars, l)
	}
	parse.columns = append(parse.columns, &TableColumn{index: 0, token: ""})
	for _, token := range tokens {
		parse.columns = append(parse.columns,
			&TableColumn{
				index:     len(parse.columns),
				token:     token.Annotations[ling.Norm],
				startByte: token.StartByte, endByte: token.EndByte,
			})
	}
	parse.finalState = parse.parse(start)
	return parse, nil
}

func extract(text string) ([]*ling.Token, *Grammar, error) {
	if text = strings.TrimSpace(text); text == "" {
		return nil, nil, fmt.Errorf("text is empty")
	}
	d := ling.NewDocument(text)
	if err := NLP().Annotate(d); err != nil {
		return nil, nil, err
	}
	var ret []*ling.Token
	for _, token := range d.Tokens {
		if token.Type == ling.Space {
			continue
		}
		ret = append(ret, token)
	}
	if len(ret) == 0 {
		return nil, nil, fmt.Errorf("no tokens")
	}
	l, err := localGrammar(d)
	if err != nil {
		return nil, nil, err
	}
	return ret, l, nil
}
