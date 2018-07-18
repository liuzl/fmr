package fmr

import (
	"fmt"
)

func (g *Grammar) MatchFrames(text string) error {
	frames, rules, err := g.getCandidates(text)
	if err != nil {
		return err
	}
	fmt.Println(frames, rules)
	return nil
}

func (g *Grammar) getCandidates(text string) (
	map[RbKey]bool, map[RbKey]bool, error) {

	hits, err := g.matcher.MultiMatch(text)
	if err != nil {
		return nil, nil, err
	}
	frames := map[RbKey]bool{}
	rules := map[RbKey]bool{}
	for _, v := range hits {
		for cate, _rbKey := range v.Value {
			rbKey, ok := _rbKey.(RbKey)
			if !ok {
				return nil, nil, fmt.Errorf("type error in grammar dict matcher")
			}
			switch cate {
			case "frame":
				frames[rbKey] = true
			case "rule":
				rules[rbKey] = true
			}
		}
	}
	var ruleList []RbKey
	for k, _ := range rules {
		ruleList = append(ruleList, k)
	}
	for {
		if len(ruleList) == 0 {
			break
		}
		r := ruleList[0]
		ruleList = ruleList[1:]
		ret, err := g.kv.Get(r.RuleName)
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
			switch cate {
			case "frame":
				frames[rbKey] = true
			case "rule":
				if !rules[rbKey] {
					ruleList = append(ruleList, rbKey)
					rules[rbKey] = true
				}
			}
		}
	}
	return frames, rules, nil
}
