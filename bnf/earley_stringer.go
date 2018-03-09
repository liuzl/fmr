package bnf

import (
	"fmt"
	"os"
)

func (ts *TableState) String() string {
	s := ""
	for i, term := range ts.rb.Terms {
		if i == ts.dot {
			s += DOT + " "
		}
		s += term.Value + " "
	}
	if ts.dot == len(ts.rb.Terms) {
		s += DOT
	}
	return fmt.Sprintf("%-6s -> %-20s [%d-%d]", ts.name, s, ts.start, ts.end)
}

func (tc *TableColumn) String() string {
	out := fmt.Sprintf("[%d] '%s'\n", tc.index, tc.token)
	out += "=======================================\n"
	for _, s := range tc.states {
		out += s.String() + "\n"
	}
	return out
}

func (p *Parser) String() string {
	out := ""
	for _, c := range p.columns {
		out += c.String() + "\n"
	}
	return out
}

func (n *Node) Print(out *os.File) {
	n.PrintLevel(out, 0)
}

func (n *Node) PrintLevel(out *os.File, level int) {
	indentation := ""
	for i := 0; i < level; i++ {
		indentation += "  "
	}
	fmt.Fprintf(out, "%s%v\n", indentation, n.value)
	for _, child := range n.children {
		child.PrintLevel(out, level+1)
	}
}

func (n *Node) String() string {
	if len(n.children) > 0 {
		return fmt.Sprintf("%+v %+v", n.value, n.children)
	} else {
		return fmt.Sprintf("%+v", n.value)
	}
}
