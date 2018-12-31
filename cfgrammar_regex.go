package fmr

import (
	"fmt"

	"github.com/liuzl/goutil"
)

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
	h := goutil.MD5(s)[:16]
	g.Regexps[h] = s
	return &Term{Value: h, Type: Nonterminal}, nil
}
