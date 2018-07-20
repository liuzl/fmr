package fmr

import (
	"fmt"
	"strings"

	"github.com/liuzl/dict"
)

func updateIndex(index map[string]map[string]interface{},
	k string, v map[string]interface{}) error {
	if index == nil {
		return fmt.Errorf("nil grammar index")
	}
	if k == "" || v == nil {
		return fmt.Errorf("empty k or v when updateIndex")
	}
	if index[k] == nil {
		index[k] = v
	} else {
		for kk, vv := range v {
			index[k][kk] = vv
		}
	}
	return nil
}

func (g *Grammar) indexRules(rules map[string]*Rule, cate string) error {
	var err error
	for _, rule := range rules {
		for id, body := range rule.Body {
			for _, term := range body.Terms {
				v := map[string]interface{}{cate: RbKey{rule.Name, id}}
				switch term.Type {
				case Terminal:
					value := strings.TrimSpace(term.Value)
					if value == "" {
						continue
					}
					if err = g.trie.SafeUpdate([]byte(value), 1); err != nil {
						return err
					}
					if err = updateIndex(g.index, value, v); err != nil {
						return err
					}
				case Nonterminal:
					if err = updateIndex(g.ruleIndex, term.Value, v); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (g *Grammar) buildIndex() error {
	if g.Refined {
		return fmt.Errorf("should call Grammar.index before Grammar.refine")
	}
	g.trie = dict.New()
	g.index = make(map[string]map[string]interface{})
	g.ruleIndex = make(map[string]map[string]interface{})

	if err := g.indexRules(g.Frames, "frame"); err != nil {
		return err
	}
	if err := g.indexRules(g.Rules, "rule"); err != nil {
		return err
	}
	return nil
}
