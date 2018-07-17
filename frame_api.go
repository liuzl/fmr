package fmr

import (
	"fmt"
)

func (g *Grammar) MatchFrames(text string) error {
	ret, err := g.matcher.MultiMatch(text)
	if err != nil {
		return err
	}
	for k, v := range ret {
		fmt.Printf("match: %s\n", k)
		for rule, typ := range v.Value {
			fmt.Printf("\t%s[%s]\n", rule, typ)
		}
	}
	return nil
}
