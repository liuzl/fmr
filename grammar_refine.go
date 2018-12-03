package fmr

import (
	"fmt"
	"strings"

	"github.com/liuzl/ling"
	"github.com/liuzl/unidecode"
	"github.com/mitchellh/hashstructure"
)

func (g *Grammar) refine(prefix string) error {
	if g.Refined {
		return nil
	}
	var terminalRules []*Rule
	var terminals = make(map[string]string)
	var names = make(map[string]bool)
	var n int
	var name string
	for _, rule := range g.Rules {
		for _, body := range rule.Body {
			for _, term := range body.Terms {
				if term.Type != Terminal {
					continue
				}
				// if this is a terminal text inside a ruleBody
				if t, has := terminals[term.Value]; has {
					term.Value = t
				} else {
					d := ling.NewDocument(term.Value)
					if err := NLP().Annotate(d); err != nil {
						return err
					}
					tname := prefix + "_t"
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
						rb.Terms = append(rb.Terms,
							&Term{Value: token.Text, Type: Terminal, Meta: term.Meta})
						if gTokens.get(token.Text) == nil {
							gTokens.put(token.Text, token)
						}
					}
					for name, n = tname, 0; ; name, n =
						fmt.Sprintf("%s_%d", tname, n), n+1 {
						if g.Rules[name] == nil && !names[name] {
							break
						}
					}
					names[name] = true
					terminals[term.Value] = name
					hash, err := hashstructure.Hash(rb, nil)
					if err != nil {
						return err
					}
					terminalRules = append(terminalRules,
						&Rule{name, map[uint64]*RuleBody{hash: rb}})
					term.Value = name
				}
				term.Type = Nonterminal
			}
		}
	}
	for _, r := range terminalRules {
		g.Rules[r.Name] = r
	}
	g.Refined = true
	return nil
}
