package fmr

import (
	"fmt"
)

func (g *Grammar) MatchFrames(text string) error {
	frames, starts, err := g.getCandidates(text)
	if err != nil {
		return err
	}
	fmt.Println(starts)
	ps, err := g.EarleyParseAll(text, starts...)
	//ps, err := g.EarleyParseAll(text, "departure")
	if err != nil {
		return err
	}
	for _, p := range ps {
		for _, finalState := range p.finalStates {
			tag := p.Tag(finalState)
			pos := p.Boundary(finalState)
			fmt.Println(tag)
			if tag == "" || pos == nil {
				return fmt.Errorf("invalid parse")
			}

			ret, err := g.kv.Get(tag)
			fmt.Println(tag, ret)
			if err != nil {
				if err.Error() == "leveldb: not found" {
					continue
				}
				return err
			}
			for cate, _rbKey := range ret {
				if cate != "frame" {
					continue
				}
				rbKey, ok := _rbKey.(RbKey)
				if !ok {
					return fmt.Errorf("format error")
				}
				if frames[rbKey] == nil {
					frames[rbKey] = &SlotFilling{make(map[Term][]*Pos), false}
				}
				t := Term{tag, Nonterminal}
				frames[rbKey].Terms[t] = append(frames[rbKey].Terms[t], pos)
			}
		}
	}
	for k, v := range frames {
		fmt.Println(k, v)
	}
	return nil
}

func (g *Grammar) getCandidates(text string) (
	map[RbKey]*SlotFilling, []string, error) {

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
					frames[rbKey] = &SlotFilling{make(map[Term][]*Pos), false}
				}
				t := Term{word, Terminal}
				for _, hit := range v.Hits {
					frames[rbKey].Terms[t] = append(frames[rbKey].Terms[t],
						&Pos{hit.Start, hit.End})
				}
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
	var starts []string
	for k, _ := range rules {
		starts = append(starts, k)
	}
	return frames, starts, nil
}
