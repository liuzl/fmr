package bnf

import (
	"fmt"
	"strings"

	"github.com/liuzl/ling"
)

func localGrammar(text string, lnlp *ling.Pipeline) (*Grammar, error) {
	if text = strings.TrimSpace(text); text == "" {
		return nil, fmt.Errorf("text is empty")
	}
	if lnlp == nil {
		lnlp = nlp
	}
	d := ling.NewDocument(text)
	if err := lnlp.Annotate(d); err != nil {
		return nil, err
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
	if err := g.refine(); err != nil {
		return nil, err
	}
	return g, nil
}
