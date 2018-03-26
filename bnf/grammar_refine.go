package bnf

import (
	"fmt"
	"github.com/liuzl/ling"
	"github.com/liuzl/unidecode"
	"strings"
)

func (g *Grammar) refine() error {
	if g.Refined {
		return nil
	}
	nlp, err := ling.NLP(ling.Norm)
	if err != nil {
		return err
	}
	var terminalRules []*Rule
	var terminals = make(map[string]string)
	var names = make(map[string]bool)
	var n = 0
	var name string
	for _, rule := range g.Rules {
		for _, body := range rule.Body {
			for _, term := range body.Terms {
				if term.IsRule {
					continue
				}
				// if this is a terminal text inside a ruleBody
				if t, has := terminals[term.Value]; has {
					term.Value = t
				} else {
					d := ling.NewDocument(term.Value)
					err = nlp.Annotate(d)
					if err != nil {
						return err
					}
					tname := "t"
					rb := &RuleBody{}
					for _, token := range d.Tokens {
						if token.Type == ling.Space {
							continue
						}
						if token.Type != ling.Punct {
							ascii := unidecode.Unidecode(token.Text)
							ascii = strings.Join(strings.Fields(ascii), "_")
							tname += "_" + ascii
						}
						rb.Terms = append(rb.Terms, &Term{token.Text, false})
					}
					for name, n = tname, 0; ; name, n = fmt.Sprintf("%s_%d", tname, n), n+1 {
						if g.Rules[name] == nil && !names[name] {
							break
						}
					}
					names[name] = true
					terminals[term.Value] = name
					terminalRules = append(terminalRules, &Rule{name, []*RuleBody{rb}})
					term.Value = name
				}
				term.IsRule = true
			}
		}
	}
	for _, r := range terminalRules {
		g.Rules[r.Name] = r
	}
	g.Refined = true
	return nil
}
