package bnf

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestTypes(t *testing.T) {
	g := &Grammar{"g1", make(map[string]*Rule)}

	g.Rules["expr"] = &Rule{"expr", []*RuleBody{
		&RuleBody{
			[]Term{Term{"a", false}},
			"",
		},
		&RuleBody{
			[]Term{Term{"expr", true}, Term{"+", false}, Term{"expr", true}},
			"nf.math.sum($1,$3)"},
	}}

	b, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
}
