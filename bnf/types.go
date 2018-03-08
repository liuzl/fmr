package bnf

type Grammar struct {
	Name  string           `json:"name"`
	Rules map[string]*Rule `json:"rules"`
}

type Rule struct {
	Name string      `json:"-"`
	Body []*RuleBody `json:"body,omitempty"`
}

type RuleBody struct {
	Terms    []*Term `json:"terms"`
	rules    []*Term `json:"-"`
	Semantic string  `json:"semantic,omitempty"`
}

type Term struct {
	Value  string `json:"value"`
	IsRule bool   `json:"is_rule"`
}
