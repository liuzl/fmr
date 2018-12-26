package fmr

import (
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"
)

func (ts *TableState) String() string {
	s := ""
	switch ts.termType {
	case Nonterminal:
		if ts.Rb != nil {
			for i, term := range ts.Rb.Terms {
				if i == ts.dot {
					s += DOT + " "
				}
				switch term.Type {
				case Nonterminal:
					s += "<" + term.Value + "> "
				case Terminal:
					s += strconv.Quote(term.Value) + " "
				case Any:
					s += "(any) "
				case List:
					s += "(list<" + term.Value + ">) "
				}
			}
			if ts.dot == len(ts.Rb.Terms) {
				s += DOT
			}
			return fmt.Sprintf("<%s> -> %s [%d-%d] {%s}",
				ts.Name, s, ts.Start, ts.End, ts.Rb.F)
		}
	case Any:
		for i := ts.Start; i < ts.End; i++ {
			s += "# "
		}
		s += DOT + " * "
		return fmt.Sprintf("(any) -> %s [%d-%d]", s, ts.Start, ts.End)
	case List:
		for i := 0; i < ts.dot; i++ {
			s += "<" + ts.Name + "> "
		}
		s += DOT + " * "
		return fmt.Sprintf("(list<%s>) -> %s [%d-%d] {fmr.list(...)}", ts.Name, s, ts.Start, ts.End)
	}
	return fmt.Sprintf("%s [%d-%d]", strconv.Quote(ts.Name), ts.Start, ts.End)
}

func (tc *TableColumn) String() (out string) {
	if tc.index == 0 {
		out = "[0] ''\n"
	} else {
		out = fmt.Sprintf("[%d] '%s' position:[%d-%d]\n",
			tc.index, tc.token, tc.token.StartByte, tc.token.EndByte)
	}
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
func (n *Node) Print(out io.Writer) {
	n.printLevel(out, 0)
}

func (n *Node) printLevel(out io.Writer, level int) {
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

func (f *FMR) String() string {
	if f == nil {
		return "nf.I($0)"
	}
	var args []string
	invalid := "invalid_fmr"
	for _, arg := range f.Args {
		switch arg.Type {
		case "string":
			if s, ok := arg.Value.(string); ok {
				args = append(args, strconv.Quote(s))
			} else {
				return invalid
			}
		case "int":
			if i, ok := arg.Value.(*big.Int); ok {
				args = append(args, i.String())
			} else {
				return invalid
			}
		case "float":
			if f, ok := arg.Value.(*big.Float); ok {
				args = append(args, f.String())
			} else {
				return invalid
			}
		case "func":
			if fmr, ok := arg.Value.(*FMR); ok {
				args = append(args, fmr.String())
			} else {
				return invalid
			}
		case "index":
			if i, ok := arg.Value.(int); ok {
				args = append(args, fmt.Sprintf("$%d", i))
			} else {
				return invalid
			}
		default:
			return invalid
		}
	}
	return fmt.Sprintf("%s(%s)", f.Fn, strings.Join(args, ","))
}
