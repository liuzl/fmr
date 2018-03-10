package bnf

import (
	"fmt"
	"os"
	"strconv"
)

func (ts *TableState) String() string {
	s := ""
	if ts.Rb != nil {
		for i, term := range ts.Rb.Terms {
			if i == ts.dot {
				s += DOT + " "
			}
			if term.IsRule {
				s += term.Value + " "
			} else {
				s += strconv.Quote(term.Value) + " "
			}
		}
		if ts.dot == len(ts.Rb.Terms) {
			s += DOT
		}
		return fmt.Sprintf("%s -> %s [%d-%d]", ts.Name, s, ts.Start, ts.End)
	}
	return fmt.Sprintf("%s [%d-%d]", strconv.Quote(ts.Name), ts.Start, ts.End)
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
	fmt.Fprintf(out, "%s%v\n", indentation, n.Value)
	for _, child := range n.Children {
		child.PrintLevel(out, level+1)
	}
}

func (n *Node) String() string {
	if len(n.Children) > 0 {
		return fmt.Sprintf("%+v %+v", n.Value, n.Children)
	} else {
		return fmt.Sprintf("%+v", n.Value)
	}
}
