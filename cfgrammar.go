package fmr

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/liuzl/goutil"
	"github.com/mitchellh/hashstructure"
)

type parser struct {
	input   string
	pos     int
	width   int
	current *position
	info    map[int]*position
	fname   string
}

type position struct {
	row, col int
	r        string
}

func (p *position) String() string {
	return fmt.Sprintf("|row:%d, col:%d, c:%s|", p.row, p.col, strconv.Quote(p.r))
}

const eof = -1

// GrammarFromFile constructs the Context-Free Grammar from file
func GrammarFromFile(file string) (*Grammar, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return grammarFromString(string(b), file, map[string]int{file: 1})
}

func grammarFromFile(file string, files map[string]int) (*Grammar, error) {
	if files[file] >= 2 {
		return nil, nil
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return grammarFromString(string(b), file, files)
}

// GrammarFromString constructs the Contex-Free Grammar from string d with name
func GrammarFromString(d, name string) (*Grammar, error) {
	return grammarFromString(d, name, make(map[string]int))
}

func grammarFromString(d, name string, files map[string]int) (*Grammar, error) {
	if files[name] >= 2 {
		return nil, nil
	}
	p := &parser{fname: name, input: d, info: make(map[int]*position)}
	if Debug {
		fmt.Println("loading ", name, files)
	}
	g, err := p.grammar(files)
	if err != nil {
		return nil, err
	}
	files[name] += 1
	if Debug {
		fmt.Println("loaded ", name, files)
	}
	return g, nil
}

func (p *parser) posInfo() string {
	return fmt.Sprintf("%s%s", p.fname, p.current)
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
	if r := p.next(); r != expected {
		return fmt.Errorf("%s :expected %s, got %s", p.posInfo(),
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
		return "", fmt.Errorf("%s : no text", p.posInfo())
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
			return "", fmt.Errorf("%s : unterminated string", p.posInfo())
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
			case '(':
				ret = append(ret, '(')
			default:
				return "", fmt.Errorf("%s : unexpected escape string", p.posInfo())
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

func (p *parser) terminal() (flags, text string, err error) {
	if err = p.eat('"'); err != nil {
		return
	}
	p.ws()
	if p.peek() == '(' {
		p.eat('(')
		p.ws()
		if err = p.eat('?'); err != nil {
			return
		}
		p.ws()
		if flags, err = p.text(); err != nil {
			return
		}
		if err = p.eat(')'); err != nil {
			return
		}
	}
	if text, err = p.terminalText(); err != nil {
		return
	}
	err = p.eat('"')
	return
}

func (p *parser) token(begin, end rune) (name string, err error) {
	if err = p.eat(begin); err != nil {
		return
	}
	if name, err = p.text(); err != nil {
		return
	}
	err = p.eat(end)
	return
}

func (p *parser) nonterminal() (string, error) {
	return p.token('<', '>')
}

func (p *parser) frame() (string, error) {
	return p.token('[', ']')
}

func (p *parser) special() (*Term, error) {
	if err := p.eat('('); err != nil {
		return nil, err
	}
	p.ws()
	name, err := p.text()
	if err != nil {
		return nil, err
	}
	p.ws()
	switch name {
	case "any":
		return p.any()
	case "list":
		return p.list()
	default:
		return nil, fmt.Errorf(
			"%s: special rule:(%s) not supported", p.posInfo(), name)
	}
}

func (p *parser) specialMeta() (map[string]int, error) {
	p.ws()
	var err error
	var meta map[string]int
	if p.peek() == '{' {
		// contains range
		meta = make(map[string]int)
		p.eat('{')
		p.ws()
		if meta["min"], err = p.getInt(); err != nil {
			return nil, err
		}
		p.ws()
		if err = p.eat(','); err != nil {
			return nil, err
		}
		p.ws()
		if meta["max"], err = p.getInt(); err != nil {
			return nil, err
		}
		if meta["max"] < meta["min"] {
			return nil, fmt.Errorf("%s : max:%d less than min:%d",
				p.posInfo(), meta["max"], meta["min"])
		}
		p.ws()
		if err = p.eat('}'); err != nil {
			return nil, err
		}
	}
	p.ws()
	return meta, nil
}

func (p *parser) list() (*Term, error) {
	name, err := p.nonterminal()
	if err != nil {
		return nil, err
	}
	meta, err := p.specialMeta()
	if err != nil {
		return nil, err
	}
	if err = p.eat(')'); err != nil {
		return nil, err
	}
	if len(meta) > 0 {
		return &Term{Type: List, Value: name, Meta: meta}, nil
	}
	return &Term{Type: List, Value: name}, nil
}

func (p *parser) any() (*Term, error) {
	meta, err := p.specialMeta()
	if err != nil {
		return nil, err
	}
	if err = p.eat(')'); err != nil {
		return nil, err
	}
	if len(meta) > 0 {
		return &Term{Type: Any, Meta: meta}, nil
	}
	return &Term{Type: Any}, nil
}

func (p *parser) regex(g *Grammar) (*Term, error) {
	if err := p.eat('`'); err != nil {
		return nil, err
	}
	p.ws()
	var ret []rune
OUT:
	for {
		switch r := p.next(); {
		case r == '`':
			break OUT
		case r == eof:
			return nil, fmt.Errorf("%s : unterminated string", p.posInfo())
		default:
			ret = append(ret, r)
		}
	}
	if len(ret) == 0 {
		return nil, fmt.Errorf("%s : empty regexp string", p.posInfo())
	}
	s := string(ret)
	if _, err := goutil.Regexp(s); err != nil {
		return nil, fmt.Errorf("%s : `%s` is not a valid regexp", p.posInfo(), s)
	}
	md5 := goutil.MD5(s)
	g.Regexps[md5] = s
	return &Term{Value: md5, Type: Nonterminal}, nil
}

func (p *parser) term(g *Grammar) (*Term, error) {
	switch p.peek() {
	case '<':
		name, err := p.nonterminal()
		if err != nil {
			return nil, err
		}
		return &Term{Value: name, Type: Nonterminal}, nil
	case '"':
		flags, text, err := p.terminal()
		if err != nil {
			return nil, err
		}
		if flags == "" {
			return &Term{Value: text, Type: Terminal}, nil
		}
		return &Term{Value: text, Type: Terminal, Meta: flags}, nil
	case '(':
		return p.special()
	case '`':
		return p.regex(g)
	}
	return nil, fmt.Errorf("%s :invalid term char", p.posInfo())
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
	var prev rune = eof
	var r rune
	first := true
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
		return "", fmt.Errorf("%s : no funcName", p.posInfo())
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
				err = fmt.Errorf("%s : unexpected semantic args", p.posInfo())
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
		err = fmt.Errorf("%s : number expected", p.posInfo())
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
	if _, text, err = p.terminal(); err != nil {
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
				return nil, fmt.Errorf("%s : unexpected dot", p.posInfo())
			}
			hasDot = true
			ret = append(ret, r)
		} else {
			break
		}
	}
	if len(ret) == 0 {
		return nil, fmt.Errorf("%s : number expected", p.posInfo())
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

func (p *parser) ruleBody(g *Grammar) (*RuleBody, error) {
	t, err := p.term(g)
	if err != nil {
		return nil, err
	}
	terms := []*Term{t}
	if err = p.comments(); err != nil {
		return nil, err
	}
	for {
		if err = p.comments(); err != nil {
			return nil, err
		}
		if !strings.ContainsRune("<\"(`", p.peek()) {
			break
		}
		if t, err = p.term(g); err != nil {
			return nil, err
		}
		terms = append(terms, t)
		if err = p.comments(); err != nil {
			return nil, err
		}
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
		if err = p.comments(); err != nil {
			return nil, err
		}
	}
	return &RuleBody{terms, f}, nil
}

func (p *parser) ruleBodies(g *Grammar) (map[uint64]*RuleBody, error) {
	r, err := p.ruleBody(g)
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
		if err = p.comments(); err != nil {
			return nil, err
		}
		if r, err = p.ruleBody(g); err != nil {
			return nil, err
		}
		if hash, err = hashstructure.Hash(r, nil); err != nil {
			return nil, err
		}
		rules[hash] = r
	}
	return rules, nil
}

func (p *parser) rule(c rune, g *Grammar) (*Rule, error) {
	var name string
	var err error
	switch c {
	case '<':
		if name, err = p.nonterminal(); err != nil {
			return nil, err
		}
	case '[':
		if name, err = p.frame(); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%s : unexpected char", p.posInfo())
	}
	if err = p.comments(); err != nil {
		return nil, err
	}
	if err = p.eat('='); err != nil {
		return nil, err
	}
	if err = p.comments(); err != nil {
		return nil, err
	}
	body, err := p.ruleBodies(g)
	if err != nil {
		return nil, err
	}
	if err = p.eat(';'); err != nil {
		return nil, err
	}
	return &Rule{name, body}, nil
}

func (p *parser) grammar(files map[string]int) (*Grammar, error) {
	g := &Grammar{
		Name:    p.fname,
		Rules:   make(map[string]*Rule),
		Frames:  make(map[string]*Rule),
		Regexps: make(map[string]string),
	}
	for {
		if err := p.comments(); err != nil {
			return nil, err
		}
		if p.peek() != '#' {
			break
		}
		p.eat('#')
		p.ws()
		name, err := p.text()
		if err != nil {
			return nil, err
		}
		if name != "include" {
			return nil, fmt.Errorf(
				"%s: directive:(%s) not suppported", p.posInfo(), name)
		}
		p.ws()
		_, ifile, err := p.terminal()
		if err != nil {
			return nil, err
		}
		files[ifile] += 1
		ig, err := grammarFromFile(ifile, files)
		if err != nil {
			return nil, err
		}
		if ig == nil {
			continue
		}
		g.includes = append(g.includes, ig)
		g.includes = append(g.includes, ig.includes...)
	}
	for {
		if err := p.comments(); err != nil {
			return nil, err
		}

		c := p.peek()
		if !strings.ContainsRune(`<[`, c) {
			break
		}
		r, err := p.rule(c, g)
		if err != nil {
			return nil, err

		}
		rules := g.Rules
		if c == '[' {
			rules = g.Frames
		}
		if _, has := rules[r.Name]; has {
			for k, v := range r.Body {
				rules[r.Name].Body[k] = v
			}
		} else {
			rules[r.Name] = r
		}
	}
	if p.next() != eof {
		return nil, fmt.Errorf("%s : format error", p.posInfo())
	}
	if err := g.buildIndex(); err != nil {
		return nil, err
	}
	if err := g.refine("g"); err != nil {
		return nil, err
	}
	return g, nil
}
