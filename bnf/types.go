package bnf

type Grammar struct {
	Name    string           `json:"name"`
	Rules   map[string]*Rule `json:"rules"`
	Refined bool             `json:"refined"`
}

type Rule struct {
	Name string      `json:"-"`
	Body []*RuleBody `json:"body,omitempty"`
}

type RuleBody struct {
	Terms []*Term `json:"terms"`
	F     *FMR    `json:"f,omitempty"`
}

type Term struct {
	Value  string `json:"value"`
	IsRule bool   `json:"is_rule"`
}

type Arg struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type FMR struct {
	Fn   string `json:fn,omitempty`
	Args []*Arg `json:args,omitempty`
}
