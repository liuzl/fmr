package bnf

import "fmt"

var Debug = false

// AST of tree structure
type Node struct {
	Value    interface{} `json:"value"`
	Children []*Node     `json:"children,omitempty"`
}

func (p *Parser) GetTrees() []*Node {
	if Debug {
		fmt.Printf("%+v\n", p)
	}
	if p.finalState != nil {
		return p.buildTrees(p.finalState)
	}
	return nil
}

func (p *Parser) buildTrees(state *TableState) []*Node {
	return p.buildTreesHelper(
		&[]*Node{}, state, len(state.Rb.Terms)-1, state.End)
}

func (p *Parser) buildTreesHelper(children *[]*Node, state *TableState,
	termIndex, end int) []*Node {
	// begin with the last --non-terminal-- of the ruleBody of finalState
	if Debug {
		fmt.Printf("%+v termIndex:%d children:%+v, end:%d\n", state, termIndex, children, end)
	}
	var outputs []*Node
	var start = -1
	if termIndex < 0 {
		// this is the base-case for the recursion (we matched the entire rule)
		outputs = append(outputs, &Node{state, *children})
		return outputs
	} else if termIndex == 0 {
		// if this is the first rule
		start = state.Start
	}
	term := state.Rb.Terms[termIndex]

	if !term.IsRule {
		n := &TableState{term.Value, nil,
			state.Start + termIndex, state.Start + termIndex + 1, 0}
		cld := []*Node{&Node{n, nil}}
		cld = append(cld, *children...)
		for _, node := range p.buildTreesHelper(&cld, state, termIndex-1, end-1) {
			outputs = append(outputs, node)
		}
		return outputs
	}

	for _, st := range p.columns[end].states {
		if st == state {
			// this prevents an endless recursion: since the states are filled in
			// order of completion, we know that X cannot depend on state Y that
			// comes after it X in chronological order
			break
		}

		if !st.isCompleted() || st.Name != term.Value {
			// this state is out of the question -- either not completed or does not
			// match the name
			if Debug {
				fmt.Printf("\tN st:%+v, term:%+v\n", st, term)
			}
			continue
		}
		if start != -1 && st.Start != start {
			// if start isn't nil, this state must span from start to end
			if Debug {
				fmt.Printf("\tN st:%+v, term:%+v\n", st, term)
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
