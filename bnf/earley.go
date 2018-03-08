package bnf

import (
	"fmt"
	"os"
)

type TableState struct {
	name       string
	production *RuleBody
	dotIndex   int
	startCol   *TableColumn
	endCol     *TableColumn
}

func (self *TableState) isCompleted() bool {
	return self.dotIndex >= self.production.size()
}

func (self *TableState) getNextTerm() *Term {
	if self.isCompleted() {
		return nil
	}
	return self.production.get(self.dotIndex)
}

func (self *TableState) String() string {
	s := ""
	for i, term := range self.production.Terms {
		if i == self.dotIndex {
			s += DOT + " "
		}
		s += term.Value + " "
	}
	if self.dotIndex == self.production.size() {
		s += DOT
	}
	return fmt.Sprintf("%-6s -> %-20s [%d-%d]",
		self.name, s, self.startCol.index, self.endCol.index)
}

type TableColumn struct {
	token  string
	index  int
	states []*TableState
}

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
	out := fmt.Sprintf("[%d] '%s'\n", self.index, self.token)
	out += "=======================================\n"
	for _, s := range self.states {
		out += s.String() + "\n"
	}
	return out
}

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

type Parser struct {
	g          *Grammar
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

// this is the name of the special "gamma" rule added by the algorithm
// (this is unicode for 'LATIN SMALL LETTER GAMMA')
const GAMMA_RULE = "\u0263" // "\u0194"
const DOT = "\u2022"        // "\u00B7"

/*
 * the Earley algorithm's core: add gamma rule, fill up table, and check if the
 * gamma rule span from the first column to the last one. return the final gamma
 * state, or null, if the parse failed.
 */
func (self *Parser) parse(start string) *TableState {
	t := &Term{start, true}
	begin := TableState{
		name:       GAMMA_RULE,
		production: &RuleBody{[]*Term{t}, []*Term{t}, ""},
		dotIndex:   0,
		startCol:   self.columns[0],
		endCol:     self.columns[0],
	}
	self.columns[0].states = append(self.columns[0].states, &begin)

	for i, col := range self.columns {
		for j := 0; j < len(col.states); j++ {
			state := col.states[j]
			if state.isCompleted() {
				self.complete(col, state)
			} else {
				term := state.getNextTerm()
				if term.IsRule {
					self.predict(col, term)
				} else if i+1 < len(self.columns) {
					self.scan(self.columns[i+1], state, term)
				}
			}
		}
		self.handleEpsilons(col)
	}

	// find end state (return nil if not found)
	lastCol := self.columns[len(self.columns)-1]
	for _, state := range lastCol.states {
		if state.name == GAMMA_RULE && state.isCompleted() {
			return state
		}
	}
	return nil
}

func (self *Parser) scan(col *TableColumn, st *TableState, term *Term) {
	if term.Value == col.token {
		col.insert(&TableState{name: st.name, production: st.production,
			dotIndex: st.dotIndex + 1, startCol: st.startCol})
	}
}

func (self *Parser) predict(col *TableColumn, term *Term) bool {
	r := self.g.Rules[term.Value] //TODO
	changed := false
	for _, prod := range r.Body {
		st := &TableState{name: r.Name, production: prod, dotIndex: 0, startCol: col}
		st2 := col.insert(st)
		changed = changed || (st == st2)
	}
	return changed
}

// Earley complete. returns true if the table has been changed, false otherwise
func (self *Parser) complete(col *TableColumn, state *TableState) bool {
	changed := false
	for _, st := range state.startCol.states {
		term := st.getNextTerm()
		if term == nil {
			continue
		}
		if term.IsRule && term.Value == state.name {
			st1 := &TableState{name: st.name, production: st.production,
				dotIndex: st.dotIndex + 1, startCol: st.startCol}
			st2 := col.insert(st1)
			changed = changed || (st1 == st2)
		}
	}
	return changed
}

func (self *Parser) handleEpsilons(col *TableColumn) {
	changed := true
	for changed {
		changed = false
		for _, state := range col.states {
			if state.isCompleted() {
				changed = changed || self.complete(col, state)
			}
			term := state.getNextTerm()
			if term != nil && term.IsRule {
				changed = changed || self.predict(col, term)
			}
		}
	}
}

func (self *Parser) GetTrees() []*Node {
	if self.finalState != nil {
		return self.buildTrees(self.finalState)
	}
	return nil
}

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

		if !st.isCompleted() || st.name != rule.Value {
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
