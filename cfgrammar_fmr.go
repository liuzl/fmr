package fmr

import (
	"fmt"
	"math/big"
	"unicode"
)

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
