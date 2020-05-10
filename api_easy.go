package fmr

// Parse returns parse trees for rule <start> at beginning
func (g *Grammar) Parse(text string, starts ...string) ([]*Node, error) {
	return g.extract(func(text string, starts ...string) ([]*Parse, error) {
		p, err := g.EarleyParse(text, starts...)
		if err != nil {
			return nil, err
		}
		return []*Parse{p}, nil
	}, text, starts...)
}

// ParseAny returns parse trees for rule <start> at any position
func (g *Grammar) ParseAny(text string, starts ...string) ([]*Node, error) {
	return g.extract(
		func(text string, starts ...string) ([]*Parse, error) {
			p, err := g.EarleyParseAny(text, starts...)
			if err != nil {
				return nil, err
			}
			return []*Parse{p}, nil
		}, text, starts...)
}

// ExtractMaxAll extracts all parse trees in text for rule <start>
func (g *Grammar) ExtractMaxAll(
	text string, starts ...string) ([]*Node, error) {
	return g.extract(g.EarleyParseMaxAll, text, starts...)
}

// ExtractAll extracts all parse trees in text for rule <start>
func (g *Grammar) ExtractAll(text string, starts ...string) ([]*Node, error) {
	return g.extract(g.EarleyParseAll, text, starts...)
}

func (g *Grammar) extract(f func(string, ...string) ([]*Parse, error),
	text string, starts ...string) ([]*Node, error) {
	ps, err := f(text, starts...)
	if err != nil {
		return nil, err
	}
	var ret []*Node
	for _, p := range ps {
		for _, f := range p.GetFinalStates() {
			ret = append(ret, p.GetTrees(f)...)
		}
	}
	return ret, nil
}
