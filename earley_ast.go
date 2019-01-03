package fmr

import "fmt"

// Debug flag
var Debug = false

// Node is the AST of tree structure
type Node struct {
	Value    *TableState `json:"value"`
	Children []*Node     `json:"children,omitempty"`

	p *Parse
}

// GetFinalStates returns the final states of p
func (p *Parse) GetFinalStates() []*TableState {
	return p.finalStates
}

// Boundary returns the start, end position in NL for finalState
func (p *Parse) Boundary(finalState *TableState) *Pos {
	if finalState == nil {
		return nil
	}
	start := p.columns[1].token.StartByte
	end := p.columns[finalState.End].token.EndByte
	return &Pos{start, end}
}

// Tag returns the Nonterminal name of finalState
func (p *Parse) Tag(finalState *TableState) string {
	if finalState == nil {
		return ""
	}
	return finalState.Rb.Terms[0].Value
}

// GetTrees returns all possible parse results
func (p *Parse) GetTrees(finalState *TableState) []*Node {
	if Debug {
		fmt.Printf("chart:\n%+v\n", p)
		fmt.Println("finalState:\n", finalState)
	}
	if finalState != nil {
		return p.buildTrees(finalState)
	}
	return nil
}

func (p *Parse) buildTrees(state *TableState) []*Node {
	if state.Term.Type == Any {
		n := &TableState{state.Term, nil, state.Start, state.End, state.End}
		cld := []*Node{{n, nil, p}}
		return cld
	}
	if state.Term.Type == List {
		state.Rb = &RuleBody{}
		var args []*Arg
		for i := 0; i < state.Dot; i++ {
			state.Rb.Terms = append(state.Rb.Terms, &Term{state.Term.Value, Nonterminal, nil})
			args = append(args, &Arg{"index", i + 1})
		}
		state.Rb.F = &FMR{"fmr.list", args}
	}
	return p.buildTreesHelper(
		&[]*Node{}, state, len(state.Rb.Terms)-1, state.End)
}

/*
 * How it works: suppose we're trying to match [X -> Y Z W]. We go from finish
 * to start, e.g., first we'll try to match W in X.encCol. Let this matching
 * state be M1. Next we'll try to match Z in M1.startCol. Let this matching
 * state be M2. And finally, we'll try to match Y in M2.startCol, which must
 * also start at X.startCol. Let this matching state be M3.
 *
 * If we matched M1, M2 and M3, then we've found a parsing for X:
 * X->
 *    Y -> M3
 *    Z -> M2
 *    W -> M1
 */
func (p *Parse) buildTreesHelper(children *[]*Node, state *TableState,
	termIndex, end int) []*Node {
	// begin with the last --non-terminal-- of the ruleBody of finalState
	if Debug {
		fmt.Printf("debug: %+v termIndex:%d children:%+v, end:%d\n",
			state, termIndex, children, end)
	}
	var outputs []*Node
	var start = -1
	if termIndex < 0 {
		// this is the base-case for the recursion (we matched the entire rule)
		outputs = append(outputs, &Node{state, *children, p})
		return outputs
	} else if termIndex == 0 {
		// if this is the first rule
		start = state.Start
	}
	term := state.Rb.Terms[termIndex]

	if term.Type == Terminal {
		n := &TableState{term, nil,
			state.Start + termIndex, state.Start + termIndex + 1, 0}
		cld := []*Node{{n, nil, p}}
		cld = append(cld, *children...)
		for _, node := range p.buildTreesHelper(&cld, state, termIndex-1, end-1) {
			outputs = append(outputs, node)
		}
		return outputs
	}

	if Debug {
		fmt.Println("\nend:", end, "term.value:", term.Value, state)
	}
	for _, st := range p.columns[end].states {
		if st == state {
			// this prevents an endless recursion: since the states are filled in
			// order of completion, we know that X cannot depend on state Y that
			// comes after it X in chronological order
			if Debug {
				fmt.Println("st==state", st, state)
				fmt.Println(p.columns[end])
			}
			break
		}
		if !st.isCompleted() || st.Term.Value != term.Value || st.Term.Type != term.Type {
			// this state is out of the question -- either not completed or does not
			// match the name
			continue
		}
		if start != -1 && st.Start != start {
			// if start isn't nil, this state must span from start to end
			continue
		}
		if Debug {
			fmt.Printf("\tY st:%+v, term:%+v\n", st, term)
		}

		// okay, so `st` matches -- now we need to create a tree for every possible
		// sub-match
		for _, subTree := range p.buildTrees(st) {
			cld := []*Node{subTree}
			cld = append(cld, *children...)
			// now try all options
			for _, node := range p.buildTreesHelper(&cld, state,
				termIndex-1, st.Start) {
				outputs = append(outputs, node)
			}
		}
	}
	return outputs
}
