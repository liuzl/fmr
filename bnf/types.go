package bnf

type Grammar struct {
	Name  string           `json:"name"`
	Rules map[string]*Rule `json:"rules"`
}

type Rule struct {
	Name string      `json:"-"`
	Body []*RuleBody `json:"body,omitempty"`
}

func (r *Rule) size() int {
	return len(r.Body)
}

func (r *Rule) get(index int) *RuleBody {
	return r.Body[index]
}

type RuleBody struct {
	Terms    []*Term `json:"terms"`
	rules    []*Term `json:"-"`
	Semantic string  `json:"semantic,omitempty"`
}

func (rb *RuleBody) size() int {
	return len(rb.Terms)
}

func (rb *RuleBody) get(index int) *Term {
	return rb.Terms[index]
}

type Term struct {
	Value  string `json:"value"`
	IsRule bool   `json:"is_rule"`
}
