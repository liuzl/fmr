package bnf

import (
	"fmt"
	"strconv"
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
	if r == '\n' {
		p.line += 1
		p.linePos = 0
	} else {
		p.linePos += w
	}
	return r
}

func (p *parser) eat(expected rune) error {
	r := p.next()
	if r != expected {
		return fmt.Errorf("|%d col %d| :expected %s, got %s",
			p.line, p.linePos,
			strconv.Quote(string(expected)), strconv.Quote(string(r)))
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

func (p *parser) semanticFn() (f *FMR, err error) {
	if err = p.eat('{'); err != nil {
		return
	}
	p.ws()
	f = &FMR{}
	if f.Fn, err = p.funcName(); err != nil {
		return
	}
	if f.Args, err = p.funcArgs(); err != nil {
		return
	}
	p.ws()
	err = p.eat('}')
	p.ws()
	return
}

func (p *parser) funcName() (string, error) {
	var ret []rune
	first := true
	var prev, r rune = eof, eof
Loop:
	for {
		switch r = p.next(); {
		case unicode.IsLetter(r) || r == '_':
			ret = append(ret, r)
		case unicode.IsDigit(r) && !first:
			ret = append(ret, r)
		case r == '.' && prev != '.' && !first:
			ret = append(ret, r)
		default:
			p.backup()
			break Loop
		}
		first = false
		prev = r
	}
	if len(ret) == 0 {
		return "", fmt.Errorf("|%d col %d| : no funcName", p.line, p.linePos)
	}
	p.ws()
	return string(ret), nil
}

func (p *parser) funcArgs() (args []*Arg, err error) {
	if err = p.eat('('); err != nil {
		return
	}
	var r rune
	var arg *Arg
	for {
		p.ws()
		switch r = p.peek(); {
		case r == '$':
			if arg, err = p.idxArg(); err != nil {
				return
			}
		case r == '"':
			if arg, err = p.strArg(); err != nil {
				return
			}
		case unicode.IsDigit(r) || r == '-':
			neg := (r == '-')
			if arg, err = p.numArg(neg); err != nil {
				return
			}
		default:
			if arg, err = p.fArg(); err != nil {
				return
			}
		}
		args = append(args, arg)
		if r == ',' {
			continue
		} else {
			p.ws()
			r = p.next()
			if r == ',' {
				continue
			} else if r == ')' {
				break
			} else {
				err = fmt.Errorf("|%d col %d| : unexpected semantic args", p.line, p.linePos)
				return
			}
		}
	}
	return
}

func (p *parser) getNumber() (idx int, err error) {
	idx = -1
	var n uint64
	var r rune
	for r = p.next(); unicode.IsDigit(r); r = p.next() {
		if n, err = strconv.ParseUint(string(r), 10, 32); err != nil {
			return
		}
		if idx == -1 {
			idx = int(n)
		} else {
			idx = idx*10 + int(n)
		}
	}
	if idx == -1 {
		err = fmt.Errorf("|%d col %d| : number expected", p.line, p.linePos)
		return
	}
	p.backup()
	return
}

func (p *parser) idxArg() (arg *Arg, err error) {
	if err = p.eat('$'); err != nil {
		return
	}
	var idx int
	if idx, err = p.getNumber(); err != nil {
		return
	}
	arg = &Arg{"index", idx}
	return
}

func (p *parser) strArg() (*Arg, error) {
	var text string
	var err error
	if text, err = p.terminal(); err != nil {
		return nil, err
	}
	return &Arg{"string", text}, nil
}

func (p *parser) numArg(neg bool) (*Arg, error) {
	var n int
	var err error
	if n, err = p.getNumber(); err != nil {
		return nil, err
	}
	if neg {
		n = -n
	}
	return &Arg{"number", n}, nil
}

func (p *parser) fArg() (*Arg, error) {
	var f *FMR
	var err error
	if f, err = p.semanticFn(); err != nil {
		return nil, err
	}
	return &Arg{"func", f}, nil
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
	var f *FMR
	if p.peek() == '{' {
		f, err = p.semanticFn()
		if err != nil {
			return nil, err
		}
	}
	return &RuleBody{terms, f}, nil
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
	err := g.refine()
	if err != nil {
		return nil, err
	}
	return g, nil
}

// CFGrammar constructs the Contex-Free Grammar from string d
func CFGrammar(d string) (*Grammar, error) {
	p := &parser{input: d}
	return p.grammar()
}
