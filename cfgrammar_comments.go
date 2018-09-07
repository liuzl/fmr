package fmr

import (
	"fmt"
)

func (p *parser) comments() error {
	defer p.ws()
	for {
		p.ws()
		c, err := p.comment()
		if err != nil {
			return err
		}
		if len(c) == 0 {
			return nil
		}
	}
}

func (p *parser) comment() (string, error) {
	if p.next() != '/' {
		p.backup()
		return "", nil
	}
	switch r := p.peek(); {
	case r == '/':
		return p.lineComment()
	case r == '*':
		return p.multiLineComment()
	default:
		return "", fmt.Errorf("%s : invalid char %s", p.posInfo(), string(r))
	}
}

func (p *parser) lineComment() (string, error) {
	if err := p.eat('/'); err != nil {
		return "", err
	}
	ret := []rune{'/', '/'}
	for {
		r := p.next()
		if r == '\n' {
			break
		}
		ret = append(ret, r)
	}
	return string(ret), nil
}

func (p *parser) multiLineComment() (string, error) {
	if err := p.eat('*'); err != nil {
		return "", err
	}
	ret := []rune{'/', '*'}
	var prev rune
	for {
		r := p.next()
		if r == eof {
			return "", fmt.Errorf("%s : unterminated string", p.posInfo())
		}
		if prev == '*' && r == '/' {
			break
		}
		ret = append(ret, r)
		prev = r
	}
	ret = append(ret, '/')
	return string(ret), nil
}
