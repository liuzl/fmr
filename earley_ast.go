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

func (p *Parse) GetFinalStates() []*TableState {
	return p.finalStates
}

func (p *Parse) Boundary(finalState *TableState) *Pos {
	if finalState == nil {
		return nil
	}
	start := p.columns[1].startByte
	end := p.columns[finalState.End].endByte
	return &Pos{start, end}
}

func (p *Parse) Tag(finalState *TableState) string {
	if finalState == nil {
		return ""
	}
	return finalState.Rb.Terms[0].Value
}

// GetTrees returns all possible parse results
func (p *Parse) GetTrees(finalState *TableState) []*Node {
	if Debug {
		fmt.Printf("%+v\n", p)
	}
	if finalState != nil {
		return p.buildTrees(finalState)
	}
	return nil
}

func (p *Parse) buildTrees(state *TableState) []*Node {
	if state.isAny {
		n := &TableState{"any", nil, state.Start, state.End, state.End, true}
		cld := []*Node{&Node{n, nil, p}}
		return cld
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

	//if !term.IsRule {
	if term.Type == Terminal {
		n := &TableState{term.Value, nil,
			state.Start + termIndex, state.Start + termIndex + 1, 0, false}
		cld := []*Node{&Node{n, nil, p}}
		cld = append(cld, *children...)
		for _, node := range p.buildTreesHelper(&cld, state, termIndex-1, end-1) {
			outputs = append(outputs, node)
		}
		return outputs
	}

	value := term.Value
	if term.Type == Any {
		value = "any"
	}

	if Debug {
		fmt.Println("\nend:", end, "term.value:", value, state)
	}
	for _, st := range p.columns[end].states {
		if st == state {
			// this prevents an endless recursion: since the states are filled in
			// order of completion, we know that X cannot depend on state Y that
			// comes after it X in chronological order
			break
		}

		if !st.isCompleted() || st.Name != value {
			// this state is out of the question -- either not completed or does not
			// match the name
			if Debug {
				//fmt.Printf("\tN st:%+v, term:%+v\n", st, term)
			}
			continue
		}
		if start != -1 && st.Start != start {
			// if start isn't nil, this state must span from start to end
			if Debug {
				//fmt.Printf("\tN st:%+v, term:%+v\n", st, term)
			}
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
			for _, node := range p.buildTreesHelper(&cld, state, termIndex-1, st.Start) {
				outputs = append(outputs, node)
			}
		}
	}
	return outputs
}
