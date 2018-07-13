package fmr

// A Grammar stores a Context-Free Grammar
type Grammar struct {
	Name    string           `json:"name"`
	Rules   map[string]*Rule `json:"rules"`
	Refined bool             `json:"refined"`
}

// A Rule stores a set of production rules of Name
type Rule struct {
	Name string               `json:"-"`
	Body map[uint64]*RuleBody `json:"body,omitempty"`
	//Body []*RuleBody `json:"body,omitempty"`
}

// A RuleBody is one production rule
type RuleBody struct {
	Terms []*Term `json:"terms"`
	F     *FMR    `json:"f,omitempty"`
}

type TermType byte

const (
	EOF TermType = iota
	Nonterminal
	Terminal
	Any
)

// A Term is the component of RuleBody
type Term struct {
	Value string   `json:"value"`
	Type  TermType `json:"type"`
}

// Arg is the type of argument for functions
type Arg struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// FMR stands for Funtional Meaning Representation
type FMR struct {
	Fn   string `json:"fn,omitempty"`
	Args []*Arg `json:"args,omitempty"`
}
