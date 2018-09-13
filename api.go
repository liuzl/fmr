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

// EarleyParse parses text for rule <start> at beginning
func (g *Grammar) EarleyParse(text string, starts ...string) (*Parse, error) {
	tokens, l, err := extract(text)
	if err != nil {
		return nil, err
	}
	return g.earleyParse(true, text, tokens, l, starts...)
}

// EarleyParseAny parses text for rule <start> at any position
func (g *Grammar) EarleyParseAny(text string, starts ...string) (*Parse, error) {
	tokens, l, err := extract(text)
	if err != nil {
		return nil, err
	}
	var p *Parse
	for i := 0; i < len(tokens); i++ {
		if p, err = g.earleyParse(true, text, tokens[i:], l, starts...); err != nil {
			return nil, err
		}
		if p.finalStates != nil {
			return p, nil
		}
	}
	return p, nil
}

// EarleyParseMaxAll extracts all submatches in text for rule <start>
func (g *Grammar) EarleyParseMaxAll(text string, starts ...string) ([]*Parse, error) {
	tokens, l, err := extract(text)
	if err != nil {
		return nil, err
	}
	var ret []*Parse
	for i := 0; i < len(tokens); {
		p, err := g.earleyParse(true, text, tokens[i:], l, starts...)
		if err != nil {
			return nil, err
		}
		if p.finalStates != nil {
			ret = append(ret, p)
			max := 0
			for _, finalState := range p.finalStates {
				if finalState.End > max {
					max = finalState.End
				}
			}
			i += max
		} else {
			i++
		}
	}
	return ret, nil
}

// EarleyParseAll extracts all submatches in text for rule <start>
func (g *Grammar) EarleyParseAll(text string, starts ...string) ([]*Parse, error) {
	tokens, l, err := extract(text)
	if err != nil {
		return nil, err
	}
	var ret []*Parse
	for i := 0; i < len(tokens); i++ {
		p, err := g.earleyParse(false, text, tokens[i:], l, starts...)
		if err != nil {
			return nil, err
		}
		if p.finalStates != nil {
			ret = append(ret, p)
			//i += p.finalState.End
		}
	}
	return ret, nil
}

func (g *Grammar) earleyParse(maxFlag bool, text string,
	tokens []*ling.Token, l *Grammar, starts ...string) (*Parse, error) {
	if len(starts) == 0 {
		return nil, fmt.Errorf("no start rules")
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens to parse")
	}

	parse := &Parse{grammars: []*Grammar{g}, text: text, starts: starts}
	if len(g.includes) > 0 {
		parse.grammars = append(parse.grammars, g.includes...)
	}
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
	parse.parse(maxFlag)
	if Debug {
		fmt.Println(parse)
	}
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
