package fmr

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mitchellh/hashstructure"
)

type parser struct {
	input   string
	pos     int
	width   int
	current *position
	info    map[int]*position
}

type position struct {
	row, col int
	r        string
}

func (p *position) String() string {
	return fmt.Sprintf("|row:%d, col:%d, c:%s|", p.row, p.col, strconv.Quote(p.r))
}

const eof = -1

// CFGrammar constructs the Contex-Free Grammar from string d
func CFGrammar(d string) (*Grammar, error) {
	p := &parser{input: d, info: make(map[int]*position)}
	return p.grammar()
}

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
	if p.info[p.pos] == nil {
		if p.current == nil {
			p.current = &position{1, w, string(r)}
		} else {
			if r == '\n' {
				p.current = &position{p.current.row + 1, w, string(r)}
			} else {
				p.current = &position{p.current.row, p.current.col + w, string(r)}
			}
		}
		p.info[p.pos] = p.current
	} else {
		p.current = p.info[p.pos]
	}
	return r
}

func (p *parser) eat(expected rune) error {
	r := p.next()
	if r != expected {
		return fmt.Errorf("%s :expected %s, got %s", p.current,
			strconv.Quote(string(expected)), strconv.Quote(string(r)))
	}
	return nil
}

func (p *parser) backup() {
	p.pos -= p.width
	p.current = p.info[p.pos]
}

func (p *parser) peek() rune {
	r := p.next()
	p.backup()
	return r
}

func (p *parser) ws() string {
	var ret []rune
	for r := p.next(); unicode.IsSpace(r); r = p.next() {
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
		return "", fmt.Errorf("%s : no text", p.current)
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
			return "", fmt.Errorf("%s : unterminated string", p.current)
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
				return "", fmt.Errorf("%s : unexpected escape string", p.current)
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
	if err = p.eat('"'); err != nil {
		return
	}
	if text, err = p.terminalText(); err != nil {
		return
	}
	err = p.eat('"')
	return
}

func (p *parser) nonterminal() (name string, err error) {
	if err = p.eat('<'); err != nil {
		return
	}
	if name, err = p.text(); err != nil {
		return
	}
	err = p.eat('>')
	return
}

func (p *parser) any() (*Term, error) {
	if err := p.eat('('); err != nil {
		return nil, err
	}
	name, err := p.text()
	if err != nil {
		return nil, err
	}
	if name != "any" {
		return nil, fmt.Errorf("%s: any rule:(%s) not supported", p.current, name)
	}
	if err := p.eat(')'); err != nil {
		return nil, err
	}
	return &Term{Type: Any}, nil
}

func (p *parser) term() (*Term, error) {
	switch p.peek() {
	case '<':
		name, err := p.nonterminal()
		if err != nil {
			return nil, err
		}
		return &Term{Value: name, Type: Nonterminal}, nil
	case '"':
		text, err := p.terminal()
		if err != nil {
			return nil, err
		}
		return &Term{Value: text, Type: Terminal}, nil
	case '(':
		return p.any()
	}
	return nil, fmt.Errorf("%s :invalid term char", p.current)
}

func (p *parser) semanticFn() (f *FMR, err error) {
	p.ws()
	f = &FMR{}
	if f.Fn, err = p.funcName(); err != nil {
		return
	}
	if f.Args, err = p.funcArgs(); err != nil {
		return
	}
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
		return "", fmt.Errorf("%s : no funcName", p.current)
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
		case unicode.IsDigit(r):
			if arg, err = p.numArg(false); err != nil {
				return
			}
		case r == '-':
			if err = p.eat('-'); err != nil {
				return
			}
			if arg, err = p.numArg(true); err != nil {
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
				err = fmt.Errorf("%s : unexpected semantic args", p.current)
				return
			}
		}
	}
	return
}

func (p *parser) getInt() (idx int, err error) {
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
		err = fmt.Errorf("%s : number expected", p.current)
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
	if idx, err = p.getInt(); err != nil {
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
	var ret []rune
	hasDot := false
	for r := p.next(); ; r = p.next() {
		if unicode.IsDigit(r) {
			ret = append(ret, r)
		} else if r == '.' {
			if hasDot {
				return nil, fmt.Errorf("%s : unexpected dot", p.current)
			}
			hasDot = true
			ret = append(ret, r)
		} else {
			break
		}
	}
	if len(ret) == 0 {
		return nil, fmt.Errorf("%s : number expected", p.current)
	}
	p.backup()
	if neg {
		ret = append([]rune{'-'}, ret...)
	}
	if hasDot {
		n := new(big.Float)
		if _, err := fmt.Sscan(string(ret), n); err != nil {
			return nil, err
		}
		return &Arg{"float", n}, nil
	}
	n := new(big.Int)
	if _, err := fmt.Sscan(string(ret), n); err != nil {
		return nil, err
	}
	return &Arg{"int", n}, nil
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
	for p.ws(); strings.ContainsRune(`<"(`, p.peek()); p.ws() {
		if t, err = p.term(); err != nil {
			return nil, err
		}
		terms = append(terms, t)
	}
	var f *FMR
	if p.peek() == '{' {
		p.eat('{')
		if f, err = p.semanticFn(); err != nil {
			return nil, err
		}
		if err = p.eat('}'); err != nil {
			return nil, err
		}
		p.ws()
	}
	return &RuleBody{terms, f}, nil
}

func (p *parser) ruleBodies() (map[uint64]*RuleBody, error) {
	r, err := p.ruleBody()
	if err != nil {
		return nil, err
	}
	hash, err := hashstructure.Hash(r, nil)
	if err != nil {
		return nil, err
	}
	rules := map[uint64]*RuleBody{hash: r}
	for {
		if p.peek() != '|' {
			break
		}
		p.eat('|')
		p.ws()
		if r, err = p.ruleBody(); err != nil {
			return nil, err
		}
		if hash, err = hashstructure.Hash(r, nil); err != nil {
			return nil, err
		}
		rules[hash] = r
	}
	return rules, nil
}

func (p *parser) rule() (*Rule, error) {
	name, err := p.nonterminal()
	if err != nil {
		return nil, err
	}
	p.ws()
	if err = p.eat('='); err != nil {
		return nil, err
	}
	p.ws()
	body, err := p.ruleBodies()
	if err != nil {
		return nil, err
	}
	if err = p.eat(';'); err != nil {
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
			//g.Rules[r.Name].Body = append(g.Rules[r.Name].Body, r.Body...)
			for k, v := range r.Body {
				g.Rules[r.Name].Body[k] = v
			}
		} else {
			g.Rules[r.Name] = r
		}
	}
	if p.next() != eof {
		return nil, fmt.Errorf("%s : format error", p.current)
	}
	if err := g.refine("g"); err != nil {
		return nil, err
	}
	return g, nil
}
