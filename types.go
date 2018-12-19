package fmr

import (
	"encoding/gob"
	"fmt"

	"github.com/liuzl/dict"
	"github.com/mitchellh/hashstructure"
)

func init() {
	gob.Register(RbKey{})
}

// A Grammar stores a Context-Free Grammar
type Grammar struct {
	Name    string            `json:"name"`
	Rules   map[string]*Rule  `json:"rules"`
	Frames  map[string]*Rule  `json:"frames"`
	Regexps map[string]string `json:"regexps"`
	Refined bool              `json:"refined"`

	trie      *dict.Cedar
	index     map[string]*Index
	ruleIndex map[string]*Index

	includes []*Grammar
}

// An Index contains two sets for frames' names and rules' names
type Index struct {
	Frames map[RbKey]struct{}
	Rules  map[RbKey]struct{}
}

// A RbKey identifies a specific RuleBody by name and id
type RbKey struct {
	RuleName string `json:"rule_name"`
	BodyId   uint64 `json:"body_id"`
}

// A Pos specifies the start and end positions
type Pos struct {
	StartByte int `json:"start_byte"`
	EndByte   int `json:"end_byte"`
}

// A Slot contains the Pos and its corresponding parse trees
type Slot struct {
	Pos
	Trees []*Node
}

// A SlotFilling is a frame consists of Slots
type SlotFilling struct {
	Fillings map[uint64][]*Slot
	Complete bool
}

func (s *SlotFilling) String() string {
	return fmt.Sprintf("Complete:%+v, %+v", s.Complete, s.Fillings)
}

// A Rule stores a set of production rules of Name
type Rule struct {
	Name string               `json:"-"`
	Body map[uint64]*RuleBody `json:"body,omitempty"`
}

// A RuleBody is one production rule
type RuleBody struct {
	Terms []*Term `json:"terms"`
	F     *FMR    `json:"f,omitempty"`
}

// TermType of grammar terms
type TermType byte

//go:generate jsonenums -type=TermType

// definition of TermTypes
const (
	EOF TermType = iota
	Nonterminal
	Terminal
	Any
	List
)

// A Term is the component of RuleBody
type Term struct {
	Value string      `json:"value"`
	Type  TermType    `json:"type"`
	Meta  interface{} `json:"meta"`
}

// Key returns a unique key for Term t
func (t *Term) Key() uint64 {
	hash, err := hashstructure.Hash(t, nil)
	if err != nil {
		return 0
	}
	return hash
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
