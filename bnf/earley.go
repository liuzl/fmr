package bnf

// this is the name of the special "gamma" rule added by the algorithm
// (this is unicode for 'LATIN SMALL LETTER GAMMA')
const GAMMA_RULE = "\u0263" // "\u0194"
const DOT = "\u2022"        // "\u00B7"

type TableState struct {
	Name  string    `json:"name"`
	Rb    *RuleBody `json:"rb,omitempty"`
	Start int       `json:"start"`
	End   int       `json:"end"`
	dot   int
}

type TableColumn struct {
	token  string
	index  int
	states []*TableState
}

type Parser struct {
	g          *Grammar
	columns    []*TableColumn
	finalState *TableState
}

func (s *TableState) isCompleted() bool {
	return s.dot >= len(s.Rb.Terms)
}

func (s *TableState) getNextTerm() *Term {
	if s.isCompleted() {
		return nil
	}
	return s.Rb.Terms[s.dot]
}

func (col *TableColumn) insert(state *TableState) *TableState {
	state.End = col.index
	for _, s := range col.states {
		if *state == *s {
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
func (p *Parser) parse(start string) *TableState {
	rb := &RuleBody{[]*Term{&Term{start, true}}, ""}
	begin := &TableState{GAMMA_RULE, rb, 0, 0, 0}
	p.columns[0].states = append(p.columns[0].states, begin)
	for i, col := range p.columns {
		for j := 0; j < len(col.states); j++ {
			state := col.states[j]
			if state.isCompleted() {
				p.complete(col, state)
			} else {
				term := state.getNextTerm()
				if term.IsRule {
					p.predict(col, term)
				} else if i+1 < len(p.columns) {
					p.scan(p.columns[i+1], state, term)
				}
			}
		}
		p.handleEpsilons(col)
	}

	// find end state (return nil if not found)
	lastCol := p.columns[len(p.columns)-1]
	for _, state := range lastCol.states {
		if state.Name == GAMMA_RULE && state.isCompleted() {
			return state
		}
	}
	return nil
}

func (*Parser) scan(col *TableColumn, st *TableState, term *Term) {
	if term.Value == col.token {
		col.insert(&TableState{Name: st.Name, Rb: st.Rb,
			dot: st.dot + 1, Start: st.Start})
	}
}

func (p *Parser) predict(col *TableColumn, term *Term) bool {
	r := p.g.Rules[term.Value] //TODO
	changed := false
	for _, prod := range r.Body {
		st := &TableState{Name: r.Name, Rb: prod, dot: 0, Start: col.index}
		st2 := col.insert(st)
		changed = changed || (st == st2)
	}
	return changed
}

// Earley complete. returns true if the table has been changed, false otherwise
func (p *Parser) complete(col *TableColumn, state *TableState) bool {
	changed := false
	for _, st := range p.columns[state.Start].states {
		term := st.getNextTerm()
		if term == nil {
			continue
		}
		if term.IsRule && term.Value == state.Name {
			st1 := &TableState{Name: st.Name, Rb: st.Rb,
				dot: st.dot + 1, Start: st.Start}
			st2 := col.insert(st1)
			changed = changed || (st1 == st2)
		}
	}
	return changed
}

func (p *Parser) handleEpsilons(col *TableColumn) {
	changed := true
	for changed {
		changed = false
		for _, state := range col.states {
			if state.isCompleted() {
				changed = changed || p.complete(col, state)
			}
			term := state.getNextTerm()
			if term != nil && term.IsRule {
				changed = changed || p.predict(col, term)
			}
		}
	}
}
