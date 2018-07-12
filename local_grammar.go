package fmr

import (
	"fmt"

	"github.com/liuzl/ling"
)

func localGrammar(d *ling.Document) (*Grammar, error) {
	if d == nil {
		return nil, fmt.Errorf("document is empty")
	}
	if len(d.Spans) == 0 {
		return nil, nil
	}
	g := &Grammar{Name: "local", Rules: make(map[string]*Rule)}
	for _, span := range d.Spans {
		if span.Annotations["value"] == nil {
			continue
		}
		m, ok := span.Annotations["value"].(map[string]interface{})
		if !ok {
			continue
		}
		terms := []*Term{&Term{span.String(), Terminal}}
		rb := &RuleBody{terms, nil}
		for k, _ := range m {
			if _, has := g.Rules[k]; has {
				g.Rules[k].Body = append(g.Rules[k].Body, rb)
			} else {
				g.Rules[k] = &Rule{k, []*RuleBody{rb}}
			}
		}
	}
	if len(g.Rules) == 0 {
		return nil, nil
	}
	if err := g.refine("l"); err != nil {
		return nil, err
	}
	return g, nil
}
