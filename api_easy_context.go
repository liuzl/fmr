package fmr

// ParseWithContext returns parse trees for rule <start> at beginning
func (g *Grammar) ParseWithContext(
	context, text string, starts ...string) ([]*Node, error) {
	return g.extractWithContext(
		func(context, text string, starts ...string) ([]*Parse, error) {
			p, err := g.EarleyParseWithContext(context, text, starts...)
			if err != nil {
				return nil, err
			}
			return []*Parse{p}, nil
		}, context, text, starts...)
}

// ParseAnyWithContext returns parse trees for rule <start> at any position
func (g *Grammar) ParseAnyWithContext(
	context, text string, starts ...string) ([]*Node, error) {
	return g.extractWithContext(
		func(context, text string, starts ...string) ([]*Parse, error) {
			p, err := g.EarleyParseAnyWithContext(context, text, starts...)
			if err != nil {
				return nil, err
			}
			return []*Parse{p}, nil
		}, context, text, starts...)
}

// ExtractMaxAllWithContext extracts all parse trees in text for rule <start>
func (g *Grammar) ExtractMaxAllWithContext(
	context, text string, starts ...string) ([]*Node, error) {
	return g.extractWithContext(
		g.EarleyParseMaxAllWithContext, context, text, starts...)
}

// ExtractAllWithContext extracts all parse trees in text for rule <start>
func (g *Grammar) ExtractAllWithContext(
	context, text string, starts ...string) ([]*Node, error) {
	return g.extractWithContext(
		g.EarleyParseAllWithContext, context, text, starts...)
}

func (g *Grammar) extractWithContext(
	f func(string, string, ...string) ([]*Parse, error),
	context, text string, starts ...string) ([]*Node, error) {
	ps, err := f(context, text, starts...)
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
