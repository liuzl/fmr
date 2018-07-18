package fmr

import (
	"fmt"

	"github.com/liuzl/d"
)

func (g *Grammar) MatchFrames(text string) error {
	frames, rules, err := g.getCandidates(text)
	if err != nil {
		return err
	}
	fmt.Println(frames, rules)
	for k, v := range frames {
		fmt.Println(k, v)
	}
	return nil
}

func (g *Grammar) getCandidates(text string) (
	map[RbKey]*SlotFilling, map[string]bool, error) {

	matches, err := g.matcher.MultiMatch(text)
	if err != nil {
		return nil, nil, err
	}
	frames := map[RbKey]*SlotFilling{}
	rules := map[string]bool{}
	for word, v := range matches {
		for cate, _rbKey := range v.Value {
			rbKey, ok := _rbKey.(RbKey)
			if !ok {
				return nil, nil, fmt.Errorf("type error in grammar dict matcher")
			}
			switch cate {
			case "frame":
				if frames[rbKey] == nil {
					frames[rbKey] = &SlotFilling{make(map[Term][]*d.Pos), false}
				}
				t := Term{word, Terminal}
				frames[rbKey].Terms[t] = append(frames[rbKey].Terms[t], v.Hits...)
			case "rule":
				rules[rbKey.RuleName] = true
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
			return nil, nil, err
		}
		for cate, _rbKey := range ret {
			rbKey, ok := _rbKey.(RbKey)
			if !ok {
				return nil, nil, fmt.Errorf("type error in grammar dicts")
			}
			if cate == "rule" && !rules[rbKey.RuleName] {
				ruleList = append(ruleList, rbKey.RuleName)
				rules[rbKey.RuleName] = true
			}
		}
	}
	return frames, rules, nil
}
