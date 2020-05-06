package fmr

import (
	"fmt"
	"strings"

	"github.com/liuzl/dict"
)

func updateIndex(index map[string]*Index, k string, cate string, v RbKey) error {
	if index == nil {
		return fmt.Errorf("nil grammar index")
	}
	if cate != "frame" && cate != "rule" {
		return fmt.Errorf("invalid cate %s", cate)
	}
	if index[k] == nil {
		index[k] = &Index{make(map[RbKey]struct{}), make(map[RbKey]struct{})}
	}
	switch cate {
	case "frame":
		index[k].Frames[v] = struct{}{}
	case "rule":
		index[k].Rules[v] = struct{}{}
	}
	return nil
}

func (g *Grammar) indexRules(rules map[string]*Rule, cate string) error {
	var err error
	for _, rule := range rules {
		for id, body := range rule.Body {
			for _, term := range body.Terms {
				v := RbKey{rule.Name, id}
				value := strings.TrimSpace(term.Value)
				if value == "" {
					continue
				}
				switch term.Type {
				case Terminal:
					if err = g.trie.SafeUpdate([]byte(value), 1); err != nil {
						return err
					}
					if err = updateIndex(g.index, value, cate, v); err != nil {
						return err
					}
				case Nonterminal:
					if err = updateIndex(g.ruleIndex, value, cate, v); err != nil {
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
	g.index = make(map[string]*Index)
	g.ruleIndex = make(map[string]*Index)

	gs := []*Grammar{g}
	gs = append(gs, g.includes...)
	for _, ig := range gs {
		if err := g.indexRules(ig.Frames, "frame"); err != nil {
			return err
		}
		if err := g.indexRules(ig.Rules, "rule"); err != nil {
			return err
		}
	}
	return nil
}
