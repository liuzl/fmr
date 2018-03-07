package bnf

import (
	"fmt"
	"strings"
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
	if r == utf8.RuneError {
		return eof
	}
	p.width = w
	p.pos += w
	p.linePos += w
	return r
}

func (p *parser) eat(expected rune) error {
	r := p.next()
	if r != expected {
		return fmt.Errorf("|%d col %d| :expected %s, got %s",
			p.line, p.linePos, expected, r)
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
	for r := p.next(); unicode.IsSpace(r); r = p.next() {
		if r == '\n' {
			p.line += 1
			p.linePos = 0
		}
		ret = append(ret, r)
	}
	p.backup()
	return string(ret)
}

func (p *parser) text() (string, error) {
	var ret []rune
	first := true
Loop:
	for {
		switch r := p.next(); {
		case unicode.IsLetter(r) || r == '_':
			ret = append(ret, r)
		case unicode.IsDigit(r) && !first:
			ret = append(ret, r)
		default:
			p.backup()
			break Loop
		}
		first = false
	}
	if len(ret) == 0 {
		return "", fmt.Errorf("|%d col %d| : no text", p.line, p.linePos)
	}
	return string(ret), nil
}

func (p *parser) terminalText() (string, error) {
	var ret []rune
	var prev rune
	for {
		switch r := p.next(); {
		case r == '"' && prev != '\\':
			p.backup()
			return string(ret), nil
		case r == eof:
			return "", fmt.Errorf("|%d col %d| : unterminated string", p.line, p.linePos)
		case prev == '\\':
			switch r {
			case '\\':
				ret = append(ret, '\\')
			case 'n':
				ret = append(ret, '\n')
			case 't':
				ret = append(ret, '\t')
			case '"':
				ret = append(ret, '"')
			default:
				return "", fmt.Errorf("|%d col %d| : unexpected escape string", p.line, p.linePos)
			}
			prev = 0
		case r == '\\':
			prev = r
		default:
			ret = append(ret, r)
			prev = r
		}
	}
	return "", fmt.Errorf("|%d col %d| : unexpected string", p.line, p.linePos)
}

func (p *parser) terminal() (text string, err error) {
	err = p.eat('"')
	if err != nil {
		return
	}
	text, err = p.terminalText()
	if err != nil {
		return
	}
	err = p.eat('"')
	return
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

func (p *parser) term() (*Term, error) {
	if p.peek() == '<' {
		name, err := p.nonterminal()
		if err != nil {
			return nil, err
		}
		return &Term{name, true}, nil
	}
	text, err := p.terminal()
	if err != nil {
		return nil, err
	}
	return &Term{text, false}, nil
}

func (p *parser) ruleBody() (*RuleBody, error) {
	t, err := p.term()
	if err != nil {
		return nil, err
	}
	terms := []*Term{t}
	p.ws()
	for p.ws(); strings.ContainsRune(`<"`, p.peek()); p.ws() {
		t, err = p.term()
		if err != nil {
			return nil, err
		}
		terms = append(terms, t)
	}
	return &RuleBody{Terms: terms}, nil
}

func (p *parser) ruleBodies() ([]*RuleBody, error) {
	r, err := p.ruleBody()
	if err != nil {
		return nil, err
	}
	rules := []*RuleBody{r}
	for {
		if p.peek() != '|' {
			break
		}
		p.eat('|')
		p.ws()
		r, err = p.ruleBody()
		if err != nil {
			return nil, err
		}
		rules = append(rules, r)
	}
	return rules, nil
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
	body, err := p.ruleBodies()
	if err != nil {
		return nil, err
	}
	err = p.eat(';')
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
	if p.next() != eof {
		return nil, fmt.Errorf("|%d col %d| : format error", p.line, p.linePos)
	}
	return g, nil
}

func Parse(d string) (*Grammar, error) {
	p := &parser{input: d}
	return p.grammar()
}
