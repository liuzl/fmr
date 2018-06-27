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
			switch term.Type {
			case Nonterminal:
				s += term.Value + " "
			case Terminal:
				s += strconv.Quote(term.Value) + " "
			case Any:
				s += "(any) "
			}
		}
		if ts.dot == len(ts.Rb.Terms) {
			s += DOT
		}
		return fmt.Sprintf("%s -> %s [%d-%d]", ts.Name, s, ts.Start, ts.End)
	}
	if ts.isAny {
		for i := ts.Start; i <= ts.End; i++ {
			if i == ts.dot+ts.Start {
				s += DOT + " "
			}
			s += "# "
		}
		return fmt.Sprintf("(any) -> %s [%d-%d]", s, ts.Start, ts.End)
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

func (p *Parse) String() string {
	out := ""
	for _, c := range p.columns {
		out += c.String() + "\n"
	}
	return out
}

// Print this tree to out
func (n *Node) Print(out *os.File) {
	n.printLevel(out, 0)
}

func (n *Node) printLevel(out *os.File, level int) {
	indentation := ""
	for i := 0; i < level; i++ {
		indentation += "  "
	}
	fmt.Fprintf(out, "%s%v\n", indentation, n.Value)
	for _, child := range n.Children {
		child.printLevel(out, level+1)
	}
}

func (n *Node) String() string {
	if len(n.Children) > 0 {
		return fmt.Sprintf("%+v %+v", n.Value, n.Children)
	}
	return fmt.Sprintf("%+v", n.Value)
}
