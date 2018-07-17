package fmr

import (
	"fmt"

	"github.com/liuzl/d"
)

func (g *Grammar) indexRules(rules map[string]*Rule, cate string) error {
	var err error
	for _, rule := range rules {
		v := map[string]interface{}{rule.Name: cate}
		for _, body := range rule.Body {
			for _, term := range body.Terms {
				fmt.Println(term, v)
				switch term.Type {
				case Terminal:
					if err = g.matcher.Update(term.Value, v); err != nil {
						return err
					}
				case Nonterminal:
					if err = g.kv.Update(term.Value, v); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (g *Grammar) index() error {
	if g.Refined {
		return fmt.Errorf("should call Grammar.index before Grammar.refine")
	}
	var err error
	if g.matcher, err = d.Load("g_matcher"); err != nil {
		return err
	}
	if g.kv, err = d.Load("g_kv"); err != nil {
		return err
	}
	if err = g.indexRules(g.Frames, "frame"); err != nil {
		return err
	}
	if err = g.indexRules(g.Rules, "rule"); err != nil {
		return err
	}
	if err = g.matcher.Save(); err != nil {
		return err
	}
	return g.kv.Save()
}
