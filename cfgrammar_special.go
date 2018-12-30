package fmr

import "fmt"

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
		return &Term{Value: "any", Type: Any, Meta: meta}, nil
	}
	return &Term{Value: "any", Type: Any}, nil
}
