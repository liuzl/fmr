package bnf

// AST of tree structure
type Node struct {
	value    interface{}
	children []*Node
}

func (p *Parser) GetTrees() []*Node {
	//fmt.Printf("%+v\n", p)
	if p.finalState != nil {
		return p.buildTrees(p.finalState)
	}
	return nil
}

func (p *Parser) buildTrees(state *TableState) []*Node {
	// build tree for state
	tree := p.buildTreesHelper(
		&[]*Node{}, state, len(state.rb.Terms)-1, state.end)
	return tree
}

func (p *Parser) buildTreesHelper(children *[]*Node, state *TableState,
	termIndex int, end int) []*Node {
	// begin with the last --non-terminal-- of the ruleBody of finalState
	var outputs []*Node
	var start = -1
	if termIndex < 0 {
		// this is the base-case for the recursion (we matched the entire rule)
		outputs = append(outputs, &Node{value: state, children: *children})
		return outputs
	} else if termIndex == 0 {
		// if this is the first rule
		start = state.start
	}
	rule := state.rb.Terms[termIndex]

	if !rule.IsRule {
		cld := []*Node{&Node{"terminal:" + rule.Value, nil}}
		cld = append(cld, *children...)
		for _, node := range p.buildTreesHelper(&cld, state, termIndex-1, end) {
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

		if !st.isCompleted() || st.name != rule.Value {
			// this state is out of the question -- either not completed or does not
			// match the name
			continue
		}
		if start != -1 && st.start != start {
			// if start isn't nil, this state must span from start to end
			continue
		}
		// okay, so `st` matches -- now we need to create a tree for every possible
		// sub-match
		for _, subTree := range p.buildTrees(st) {
			// in python: children2 = [subTree] + children
			children2 := []*Node{subTree}
			children2 = append(children2, *children...)
			// now try all options
			for _, node := range p.buildTreesHelper(
				&children2, state, termIndex-1, st.start) {
				outputs = append(outputs, node)
			}
		}
	}
	return outputs
}
