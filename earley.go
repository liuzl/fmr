package fmr

import (
	"fmt"
)

// GAMMA_RULE is the name of the special "gamma" rule added by the algorithm
// (this is unicode for 'LATIN SMALL LETTER GAMMA')
const GAMMA_RULE = "\u0263" // "\u0194"

// DOT indicates the current position inside a TableState
const DOT = "\u2022" // "\u00B7"

// TableState uses Earley's dot notation: given a production X → αβ,
// the notation X → α • β represents a condition in which α has already
// been parsed and β is expected.
type TableState struct {
	Name  string    `json:"name"`
	Rb    *RuleBody `json:"rb,omitempty"`
	Start int       `json:"start"`
	End   int       `json:"end"`
	dot   int
	isAny bool
	meta  interface{}
}

// Equal func for TableState
func (s *TableState) Equal(ts *TableState) bool {
	if s == nil && ts == nil {
		return true
	}
	if !(s != nil && ts != nil) {
		if Debug {
			fmt.Println("only one is nil:", s, ts)
		}
		return false
	}
	if s.Name == ts.Name && s.Rb.Equal(ts.Rb) &&
		s.Start == ts.Start && s.End == ts.End &&
		s.dot == ts.dot && s.isAny == ts.isAny {
		return metaEqual(s.meta, ts.meta)
	}
	return false
}

func (s *TableState) metaEmpty() bool {
	if s.meta == nil {
		return true
	}
	if m, ok := s.meta.(map[string]int); ok && len(m) == 0 {
		return true
	}
	return false
}

// TableColumn is the TableState set
type TableColumn struct {
	token     string
	startByte int
	endByte   int
	index     int
	states    []*TableState
}

// Parse stores a parse chart by grammars
type Parse struct {
	grammars    []*Grammar
	text        string
	starts      []string
	columns     []*TableColumn
	finalStates []*TableState
}

func (s *TableState) isCompleted() bool {
	if s.isAny {
		if s.metaEmpty() {
			if s.dot > 0 {
				return true
			}
		} else {
			if meta, ok := s.meta.(map[string]int); ok {
				if s.dot >= meta["min"] && s.dot <= meta["max"] {
					return true
				}
			}
		}
		return false
	}
	return s.dot >= len(s.Rb.Terms)
}

func (s *TableState) getNextTerm() *Term {
	if s.isCompleted() {
		return nil
	}
	if s.isAny {
		return &Term{Type: Any, Meta: s.meta}
	}
	return s.Rb.Terms[s.dot]
}

func (col *TableColumn) insert(state *TableState) *TableState {
	return col.insertToEnd(state, false)
}

func (col *TableColumn) insertToEnd(state *TableState, end bool) *TableState {
	state.End = col.index
	if state.isAny {
		state.dot = state.End - state.Start
	}
	for i, s := range col.states {
		if s.Equal(state) {
			if end {
				col.states = append(col.states[:i], col.states[i+1:]...)
				col.states = append(col.states, s)
			}
			return s
		}
	}
	col.states = append(col.states, state)
	return col.states[len(col.states)-1]
}

/*
 * the Earley algorithm's core: add gamma rule, fill up table, and check if the
 * gamma rule span from the first column to the last one. return the final gamma
 * state, or null, if the parse failed.
 */
func (p *Parse) parse(maxFlag bool) []*TableState {
	if len(p.starts) == 0 {
		return nil
	}
	for _, start := range p.starts {
		rb := &RuleBody{
			[]*Term{{Value: start, Type: Nonterminal}},
			&FMR{"nf.I", []*Arg{{"index", 1}}},
		}
		begin := &TableState{GAMMA_RULE, rb, 0, 0, 0, false, nil}
		p.columns[0].states = append(p.columns[0].states, begin)
	}
	for i, col := range p.columns {
		if Debug {
			fmt.Printf("Column %d[%s]:\n", i, col.token)
		}
		for j := 0; j < len(col.states); j++ {
			st := col.states[j]
			if Debug {
				fmt.Printf("\tRow %d: %+v, len:%d\n", j, st, len(col.states))
			}
			if st.isAny {
				if st.metaEmpty() {
					if st.dot > 0 {
						p.complete(col, st)
					}
					if i+1 < len(p.columns) {
						p.scan(p.columns[i+1], st, &Term{Type: Any})
					}
				} else {
					if meta, ok := st.meta.(map[string]int); ok {
						if st.dot >= meta["min"] && st.dot <= meta["max"] {
							p.complete(col, st)
						}
						if i+1 < len(p.columns) && st.dot+1 <= meta["max"] {
							p.scan(p.columns[i+1], st,
								&Term{Type: Any, Meta: st.meta})
						}
					}
				}
			} else {
				if st.isCompleted() {
					p.complete(col, st)
				} else {
					term := st.getNextTerm()
					switch term.Type {
					case Nonterminal, Any:
						p.predict(col, term)
					case Terminal:
						if i+1 < len(p.columns) {
							p.scan(p.columns[i+1], st, term)
						}
					}
				}
			}
			if Debug {
				fmt.Println()
			}
		}
		if Debug {
			fmt.Println()
		}
		//p.handleEpsilons(col)
	}

	// find end state (return nil if not found)
	/*
		lastCol := p.columns[len(p.columns)-1]
		for _, state := range lastCol.states {
			if state.Name == GAMMA_RULE && state.isCompleted() {
				return state
			}
		}
	*/
	var ret []*TableState
	for i := len(p.columns) - 1; i >= 0; i-- {
		for _, state := range p.columns[i].states {
			if state.Name == GAMMA_RULE && state.isCompleted() {
				ret = append(ret, state)
				if maxFlag {
					p.finalStates = ret
					return ret
				}
			}
		}
	}
	p.finalStates = ret
	return ret
}

func (*Parse) scan(col *TableColumn, st *TableState, term *Term) {
	if term.Type == Any {
		newSt := &TableState{Name: "any", Rb: st.Rb,
			dot: st.dot + 1, Start: st.Start, isAny: st.isAny, meta: term.Meta}
		col.insert(newSt)
		if Debug {
			fmt.Println("\tscan Any")
			fmt.Printf("\t\tinsert to next: %+v\n", newSt)
		}
		return
	}
	if term.Value == col.token {
		newSt := &TableState{Name: st.Name, Rb: st.Rb,
			dot: st.dot + 1, Start: st.Start}
		col.insert(newSt)
		if Debug {
			fmt.Println("\tscan", term.Value, col.token)
			fmt.Printf("\t\tinsert to next: %+v\n", newSt)
		}
	}
}

func predict(g *Grammar, col *TableColumn, term *Term) bool {
	r, has := g.Rules[term.Value]
	if !has {
		return false
	}
	changed := false
	for _, prod := range r.Body {
		st := &TableState{Name: r.Name, Rb: prod, dot: 0, Start: col.index}
		st2 := col.insert(st)
		if Debug {
			fmt.Printf("\t\tinsert: %+v\n", st)
		}
		changed = changed || (st == st2)
	}
	return changed
}

func (p *Parse) predict(col *TableColumn, term *Term) bool {
	if Debug {
		fmt.Println("\tpredict", term.Type, term.Value)
	}
	switch term.Type {
	case Nonterminal:
		changed := false
		for _, g := range p.grammars {
			changed = predict(g, col, term) || changed
		}
		return changed
	case Any:
		st := &TableState{
			Name: "any", Start: col.index, isAny: true, meta: term.Meta}
		st2 := col.insert(st)
		fmt.Printf("\t\tinsert: %+v\n", st)
		return st == st2
	}
	return false
}

// Earley complete. returns true if the table has been changed, false otherwise
func (p *Parse) complete(col *TableColumn, state *TableState) bool {
	if Debug {
		fmt.Println("\tcomplete")
	}
	changed := false
	for _, st := range p.columns[state.Start].states {
		next := st.getNextTerm()
		if next == nil {
			continue
		}
		if (next.Type == Any && state.isAny) ||
			(next.Type == Nonterminal && next.Value == state.Name) {
			st1 := &TableState{Name: st.Name, Rb: st.Rb,
				dot: st.dot + 1, Start: st.Start, isAny: st.isAny, meta: next.Meta}
			//st2 := col.insertToEnd(st1, true)
			st2 := col.insertToEnd(st1, false)
			if Debug {
				fmt.Printf("\t\tinsert: %+v\n", st1)
			}
			changed = changed || (st1 == st2)
		}
	}
	return changed
}

func (p *Parse) handleEpsilons(col *TableColumn) {
	changed := true
	for changed {
		changed = false
		for _, state := range col.states {
			if state.isCompleted() {
				changed = p.complete(col, state) || changed
			}
			term := state.getNextTerm()
			if term != nil && term.Type == Nonterminal {
				changed = p.predict(col, term) || changed
			}
		}
	}
}
