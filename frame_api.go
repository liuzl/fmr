package fmr

import (
	"fmt"
)

func (g *Grammar) MatchFrames(text string) error {
	hits, err := g.matcher.MultiMatch(text)
	if err != nil {
		return err
	}
	frames := map[string]bool{}
	rules := map[string]bool{}
	for _, v := range hits {
		for cate, _rule := range v.Value {
			rule, ok := _rule.(string)
			if !ok {
				return fmt.Errorf("type error in grammar dict matcher")
			}
			switch cate {
			case "frame":
				frames[rule] = true
			case "rule":
				rules[rule] = true
			}
		}
	}
	var ruleList []string
	for k, _ := range rules {
		ruleList = append(ruleList, k)
	}
	for {
		if len(ruleList) == 0 {
			break
		}
		r := ruleList[0]
		ruleList = ruleList[1:]
		ret, err := g.kv.Get(r)
		if err != nil {
			if err.Error() == "leveldb: not found" {
				continue
			}
			return err
		}
		for cate, _rule := range ret {
			rule, ok := _rule.(string)
			if !ok {
				return fmt.Errorf("type error in grammar dicts")
			}
			switch cate {
			case "frame":
				frames[rule] = true
			case "rule":
				if !rules[rule] {
					ruleList = append(ruleList, rule)
					rules[rule] = true
				}
			}
		}
	}
	fmt.Println(frames, rules)
	return nil
}
