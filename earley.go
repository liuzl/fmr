package earley

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

/*
 * Terminology
 * ==========
 * Consider the following context-free rule:
 *
 *     X -> A B C | A hello
 *
 * We say rule 'X' has two __production__: "A B C" and "A hello".
 * Each production is made of __production terms__, which can be either
 * __terminals__ (in our case, "hello") or __rules__ (non-terminals, such
 * as "A", "B", and "C")
 */

/*
 * Represents a terminal element in a production
 */
type Terminal struct {
	value string
}

func (self *Terminal) String() string {
	return self.value
}

/*
 * Represents a production of the rule.
 */
type Production struct {
	terms []interface{}
	rules []*Rule
}

func NewProduction(terms ...interface{}) *Production {
	prod := &Production{}
	for _, term := range terms {
		switch term.(type) {
		case *Terminal:
			prod.terms = append(prod.terms, term.(*Terminal))
		case string: // treat string as Terminal
			prod.terms = append(prod.terms, &Terminal{term.(string)})
		case *Rule:
			prod.terms = append(prod.terms, term.(*Rule))
		default:
			println(reflect.TypeOf(term).String() + " invalid type")
		}
	}
	prod.getRules()
	return prod
}

func (self *Production) size() int {
	return len(self.terms)
}

func (self *Production) get(index int) interface{} {
	return self.terms[index]
}

func (self *Production) getRules() {
	self.rules = nil
	for _, term := range self.terms {
		switch term.(type) {
		case *Rule:
			self.rules = append(self.rules, term.(*Rule))
		}
	}
}

func (self *Production) String() string {
	s := ""
	for i, term := range self.terms {
		switch term.(type) {
		case *Terminal:
			s += term.(*Terminal).value
		case *Rule:
			s += term.(*Rule).name
		}
		if i != self.size()-1 {
			s += " "
		}
	}
	return s
}

// Epsilon transition: an empty production
var Epsilon = Production{}

/*
 * A CFG rule. Since CFG rules can be self-referential, more productions may be added
 * to them after construction. For example:
 *
 * Grammar:
 *    SYM -> a
 *    OP -> + | -
 *    EXPR -> SYM | EXPR OP EXPR
 *
 * In Golang:
 *     SYM := NewRule("SYM", NewProduction(&Terminal{"a"}))
 *     OP := NewRule("OP", NewProduction("+"))
 *     EXPR := NewRule("EXPR", NewProduction(SYM))
 *     EXPR.add(NewProduction(EXPR, OP, EXPR))
 *
 */

type Rule struct {
	name        string
	productions []*Production
}

func NewRule(name string, prods ...*Production) *Rule {
	return &Rule{name: name, productions: prods}
}

func (self *Rule) add(prods ...*Production) {
	self.productions = append(self.productions, prods...)
}

func (self *Rule) size() int {
	return len(self.productions)
}

func (self *Rule) get(index int) *Production {
	return self.productions[index]
}

func (self *Rule) String() string {
	s := self.name + " -> "
	for i, prod := range self.productions {
		s += prod.String()
		if i != self.size()-1 {
			s += " | "
		}
	}
	return s
}

/*
 * Represents a state in the Earley parsing table. A state has its rule's name,
 * the rule's production, dot-location, and starting- and ending-column in the
 * parsing table
 */
type TableState struct {
	name       string
	production *Production
	dotIndex   int
	startCol   *TableColumn
	endCol     *TableColumn
}

func (self *TableState) isCompleted() bool {
	return self.dotIndex >= self.production.size()
}

func (self *TableState) getNextTerm() interface{} {
	if self.isCompleted() {
		return nil
	}
	return self.production.get(self.dotIndex)
}

func (self *TableState) String() string {
	s := ""
	for i, term := range self.production.terms {
		if i == self.dotIndex {
			s += "\u00B7"
		}
		switch term.(type) {
		case *Terminal:
			s += term.(*Terminal).value
		case *Rule:
			s += term.(*Rule).name
		}
		s += " "
	}
	if self.dotIndex == self.production.size() {
		s += "\u00B7"
	}
	return fmt.Sprintf("%s -> %s [%d-%d]",
		self.name, s, self.startCol.index, self.endCol.index)
}

/*
 * Represents a column in the Earley parsing table
 */
type TableColumn struct {
	token  string
	index  int
	states []*TableState
}

/*
 * only insert a state if it is not already contained in the list of states.
 * return the inserted state, or the pre-existing one.
 */
func (self *TableColumn) insert(state *TableState) *TableState {
	state.endCol = self
	for _, s := range self.states {
		if *state == *s {
			return s
		}
	}
	self.states = append(self.states, state)
	return self.get(self.size() - 1)
}

func (self *TableColumn) size() int {
	return len(self.states)
}

func (self *TableColumn) get(index int) *TableState {
	return self.states[index]
}

func (self *TableColumn) String() string {
	out := ""
	out += fmt.Sprintf("[%d] '%s'\n", self.index, self.token)
	out += "=======================================\n"
	for _, s := range self.states {
		out += s.String() + "\n"
	}
	return out
}

func (self *TableColumn) Print(out *os.File, showUncompleted bool) {
	fmt.Fprintf(out, "[%d] '%s'\n", self.index, self.token)
	fmt.Fprintln(out, "=======================================")
	for _, s := range self.states {
		if !s.isCompleted() && !showUncompleted {
			continue
		}
		fmt.Fprintln(out, s)
	}
	fmt.Fprintln(out)
}

/*
 * A generic tree node
 */
type Node struct {
	value    interface{}
	children []*Node
}

func (self *Node) Print(out *os.File) {
	self.PrintLevel(out, 0)
}

func (self *Node) PrintLevel(out *os.File, level int) {
	indentation := ""
	for i := 0; i < level; i++ {
		indentation += "  "
	}
	fmt.Fprintf(out, "%s%v\n", indentation, self.value)
	for _, child := range self.children {
		child.PrintLevel(out, level+1)
	}
}

func (self *Node) String() string {
	if len(self.children) > 0 {
		return fmt.Sprintf("%+v %+v", self.value, self.children)
	} else {
		return fmt.Sprintf("%+v", self.value)
	}
}

/*
 * The Earley Parser.
 *
 * Usage:
 *
 *   var p *Parser = NewParser(StartRule, "my space-delimited statement")
 *   for _, tree := range p.getTrees() {
 *     tree.Print(os.Stdout)
 *   }
 *
 */
type Parser struct {
	columns    []*TableColumn
	finalState *TableState
}

func (self *Parser) String() string {
	out := ""
	for _, c := range self.columns {
		out += c.String() + "\n"
	}
	return out
}

func NewParser(startRule *Rule, text string) *Parser {
	tokens := strings.Fields(text)
	parser := &Parser{}
	parser.columns = append(parser.columns, &TableColumn{index: 0, token: ""})
	for i, token := range tokens {
		parser.columns = append(parser.columns,
			&TableColumn{index: i + 1, token: token})
	}
	parser.finalState = parser.parse(startRule)
	return parser
}

// this is the name of the special "gamma" rule added by the algorithm
// (this is unicode for 'LATIN SMALL LETTER GAMMA')
const GAMMA_RULE = "\u0263" // "\u0194"

/*
 * the Earley algorithm's core: add gamma rule, fill up table, and check if the
 * gamma rule span from the first column to the last one. return the final gamma
 * state, or null, if the parse failed.
 */
func (self *Parser) parse(startRule *Rule) *TableState {
	begin := TableState{
		name:       GAMMA_RULE,
		production: NewProduction(startRule),
		dotIndex:   0,
		startCol:   self.columns[0],
		endCol:     self.columns[0]}
	self.columns[0].states = append(self.columns[0].states, &begin)

	for i, col := range self.columns {
		for j := 0; j < len(col.states); j++ {
			state := col.states[j]
			if state.isCompleted() {
				self.complete(col, state)
			} else {
				var term interface{} = state.getNextTerm()
				switch term.(type) {
				case *Rule:
					self.predict(col, term.(*Rule))
				case *Terminal:
					if i+1 < len(self.columns) {
						self.scan(self.columns[i+1], state, term.(*Terminal).value)
					}
				}
			}
		}
		self.handleEpsilons(col)
		// DEBUG -- uncomment to print the table during parsing, column after column
		//col.Print(os.Stdout, true)
	}

	// find end state (return nil if not found)
	lastCol := self.columns[len(self.columns)-1]
	for i := 0; i < len(lastCol.states); i++ {
		if lastCol.states[i].name == GAMMA_RULE && lastCol.states[i].isCompleted() {
			return lastCol.states[i]
		}
	}
	return nil
}

/*
 * Earley scan
 */
func (self *Parser) scan(col *TableColumn, st *TableState, token string) {
	if token == col.token {
		col.insert(&TableState{name: st.name, production: st.production,
			dotIndex: st.dotIndex + 1, startCol: st.startCol})
	}
}

/*
 * Earley predict. returns true if the table has been changed, false otherwise
 */
func (self *Parser) predict(col *TableColumn, r *Rule) bool {
	changed := false
	for _, prod := range r.productions {
		st := &TableState{name: r.name, production: prod, dotIndex: 0, startCol: col}
		st2 := col.insert(st)
		changed = changed || (st == st2)
	}
	return changed
}

/*
 * Earley complete. returns true if the table has been changed, false otherwise
 */
func (self *Parser) complete(col *TableColumn, state *TableState) bool {
	changed := false
	for _, st := range state.startCol.states {
		var term interface{} = st.getNextTerm()
		if r, ok := term.(*Rule); ok && r.name == state.name {
			st1 := &TableState{name: st.name, production: st.production,
				dotIndex: st.dotIndex + 1, startCol: st.startCol, endCol: col}
			st2 := col.insert(st1)
			changed = changed || (st1 == st2)
		}
	}
	return changed
}

/*
 * call predict() and complete() for as long as the table keeps changing (may only
 * happen if we've got epsilon transitions)
 */
func (self *Parser) handleEpsilons(col *TableColumn) {
	changed := true
	for changed {
		changed = false
		for _, state := range col.states {
			var term interface{} = state.getNextTerm()
			if r, ok := term.(*Rule); ok {
				changed = changed || self.predict(col, r)
			}
			if state.isCompleted() {
				changed = changed || self.complete(col, state)
			}
		}
	}
}

/*
 * return all parse trees (forest). the forest is simply a list of root nodes, each
 * representing a possible parse tree. a node is contains a value and the node's
 * children, and supports pretty-printing
 */
func (self *Parser) getTrees() []*Node {
	if self.finalState != nil {
		return self.buildTrees(self.finalState)
	}
	return nil
}

/*
 * how it works: suppose we're trying to match [X -> Y Z W]. we go from finish to
 * start, e.g., first we'll try to match W in X.encCol. let this matching state be
 * M1. next we'll try to match Z in M1.startCol. let this matching state be M2. and
 * finally, we'll try to match Y in M2.startCol, which must also start at X.startCol.
 * let this matching state be M3.
 *
 * if we matched M1, M2 and M3, then we've found a parsing for X:
 * X->
 *    Y -> M3
 *    Z -> M2
 *    W -> M1
 */

func (self *Parser) buildTrees(state *TableState) []*Node {
	return self.buildTreesHelper(
		&[]*Node{}, state, len(state.production.rules)-1, state.endCol)
}

func (self *Parser) buildTreesHelper(children *[]*Node, state *TableState,
	ruleIndex int, endCol *TableColumn) []*Node {
	// begin with the last --non-terminal-- of the production of finalState
	var outputs []*Node
	var startCol *TableColumn
	if ruleIndex < 0 {
		// this is the base-case for the recursion (we matched the entire rule)
		outputs = append(outputs, &Node{value: state, children: *children})
		return outputs
	} else if ruleIndex == 0 {
		// if this is the first rule
		startCol = state.startCol
	}
	rule := state.production.rules[ruleIndex]

	for _, st := range endCol.states {
		if st == state {
			// this prevents an endless recursion: since the states are filled in
			// order of completion, we know that X cannot depend on state Y that
			// comes after it X in chronological order
			break
		}
		if !st.isCompleted() || st.name != rule.name {
			// this state is out of the question -- either not completed or does not
			// match the name
			continue
		}
		if startCol != nil && st.startCol != startCol {
			// if startCol isn't nil, this state must span from startCol to endCol
			continue
		}
		// okay, so `st` matches -- now we need to create a tree for every possible
		// sub-match
		for _, subTree := range self.buildTrees(st) {
			// in python: children2 = [subTree] + children
			children2 := []*Node{}
			children2 = append(children2, subTree)
			children2 = append(children2, *children...)
			// now try all options
			for _, node := range self.buildTreesHelper(
				&children2, state, ruleIndex-1, st.startCol) {
				outputs = append(outputs, node)
			}
		}
	}
	return outputs
}
