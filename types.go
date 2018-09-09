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
	Name    string           `json:"name"`
	Rules   map[string]*Rule `json:"rules"`
	Frames  map[string]*Rule `json:"frames"`
	Refined bool             `json:"refined"`

	trie      *dict.Cedar
	index     map[string]*Index
	ruleIndex map[string]*Index

	includes []*Grammar
}

type Index struct {
	Frames map[RbKey]struct{}
	Rules  map[RbKey]struct{}
}

// A RbKey identifies a specific RuleBody by name and id
type RbKey struct {
	RuleName string `json:"rule_name"`
	BodyId   uint64 `json:"body_id"`
}

type Pos struct {
	StartByte int `json:"start_byte"`
	EndByte   int `json:"end_byte"`
}

type Slot struct {
	Pos
	Trees []*Node
}

type SlotFilling struct {
	//Fillings map[Term][]*Slot
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

type TermType byte

//go:generate jsonenums -type=TermType

const (
	EOF TermType = iota
	Nonterminal
	Terminal
	Any
)

// A Term is the component of RuleBody
type Term struct {
	Value string      `json:"value"`
	Type  TermType    `json:"type"`
	Meta  interface{} `json:"meta"`
}

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
