package cfg

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type parser struct {
	input   string
	line    int
	pos     int
	width   int
	linePos int
}

const eof = -1

func (p *parser) next() rune {
	if p.pos >= len(p.input) {
		p.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(p.input[p.pos:])
	p.width = w
	p.pos += w
	p.linePos += w
	return r
}

func (p *parser) eat(expected rune) error {
	r := p.next()
	if r != expected {
		return fmt.Errorf("Expected %s, got %s", expected, r)
	}
	return nil
}

func (p *parser) backup() {
	p.pos -= p.width
}

func (p *parser) peek() rune {
	r := p.next()
	p.backup()
	return r
}

func (p *parser) ws() string {
	var ret []rune
	for r := p.peek(); unicode.IsSpace(r); {
		if r == '\n' {
			p.line += 1
			p.linePos = 0
		}
		ret = append(ret, r)
		p.eat(r)
	}
	return string(ret)
}

func (p *parser) text() (string, error) {
	return "", nil
}

func (p *parser) nonterminal() (name string, err error) {
	err = p.eat('<')
	if err != nil {
		return
	}
	name, err = p.text()
	if err != nil {
		return
	}
	err = p.eat('>')
	return
}
func (p *parser) expressions() ([]*RuleBody, error) {
	return nil, nil
}

func (p *parser) rule() (*Rule, error) {
	name, err := p.nonterminal()
	if err != nil {
		return nil, err
	}
	p.ws()
	err = p.eat('=')
	if err != nil {
		return nil, err
	}
	p.ws()
	body, err := p.expressions()
	if err != nil {
		return nil, err
	}
	return &Rule{name, body}, nil
}

func (p *parser) grammar() (*Grammar, error) {
	g := &Grammar{Name: "grammar", Rules: make(map[string]*Rule)}
	for p.ws(); p.peek() == '<'; p.ws() {
		r, err := p.rule()
		if err != nil {
			return nil, err
		}
		if _, has := g.Rules[r.Name]; has {
			g.Rules[r.Name].Body = append(g.Rules[r.Name].Body, r.Body...)
		} else {
			g.Rules[r.Name] = r
		}
	}
	return g, nil
}

func Parse(d string) (*Grammar, error) {
	p := &parser{input: d}
	return p.grammar()
}
