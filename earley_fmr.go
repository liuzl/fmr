package fmr

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/liuzl/goutil"
)

// Pos returns the corresponding pos of Node n in original text
func (n *Node) Pos() *Pos {
	return n.p.Boundary(n.Value)
}

// Term returns the root Term of tree node
func (n *Node) Term() *Term {
	if n.Value == nil || n.Value.Rb == nil || len(n.Value.Rb.Terms) < 1 {
		return nil
	}
	//fmt.Println(n.Children[0].Value.Rb.F)
	//fmt.Println(n.Value.Name)
	if n.Value.Term.Value == GammaRule {
		return n.Value.Rb.Terms[0]
	}
	return n.Value.Term
}

// OriginalText returns the original text of Node n
func (n *Node) OriginalText() string {
	pos := n.Pos()
	return n.p.text[pos.StartByte:pos.EndByte]
}

// NL returns the normalized text of Node n
func (n *Node) NL() string {
	var s []string
	for i := n.Value.Start + 1; i <= n.Value.End; i++ {
		s = append(s, n.p.columns[i].token.Text)
	}
	return goutil.Join(s)
}

// Semantic returns the stringified FMR of Node n
func (n *Node) Semantic() (string, error) {
	nl := strconv.Quote(n.NL())
	if n.Value.Rb == nil || n.Value.Rb.F == nil {
		if n.p == nil {
			return "", nil
		}
		// by default, returns nf.I($0)
		return nl, nil
	}
	return fmrStr(n.Value.Rb.F, n.Children, nl)
}

func fmrStr(f *FMR, children []*Node, nl string) (string, error) {
	if f == nil {
		return "", nil
	}
	if f.Fn == "nf.I" {
		if len(f.Args) != 1 {
			return "", fmt.Errorf("the length of Args of nf.I should be one")
		}
		s, err := semStr(f.Args[0], children, nl)
		if err != nil {
			return "", err
		}
		return s, nil
	}

	var args []string
	for _, arg := range f.Args {
		s, err := semStr(arg, children, nl)
		if err != nil {
			return "", err
		}
		args = append(args, s)
	}
	return fmt.Sprintf("%s(%s)", f.Fn, strings.Join(args, ", ")), nil
}

func semStr(arg *Arg, nodes []*Node, nl string) (string, error) {
	if arg == nil {
		return "", fmt.Errorf("arg is nil")
	}
	switch arg.Type {
	case "string":
		if s, ok := arg.Value.(string); ok {
			return strconv.Quote(s), nil
		}
		return "", fmt.Errorf("arg.Value: %+v is not string", arg.Value)
	case "int":
		if i, ok := arg.Value.(*big.Int); ok {
			return i.String(), nil
		}
		return "", fmt.Errorf("arg.Value: %+v is not int", arg.Value)
	case "float":
		if f, ok := arg.Value.(*big.Float); ok {
			return f.String(), nil
		}
		return "", fmt.Errorf("arg.Value: %+v is not float", arg.Value)
	case "func":
		if fmr, ok := arg.Value.(*FMR); ok {
			return fmrStr(fmr, nodes, nl)
		}
		return "", fmt.Errorf("arg.Value: %+v is not func", arg.Value)
	case "index":
		i, ok := arg.Value.(int)
		if !ok {
			return "", fmt.Errorf("arg.Value: %+v is not index", arg.Value)
		}
		if i < 0 || i > len(nodes) {
			return "", fmt.Errorf("i=%d not in range [0, %d]", i, len(nodes))
		}
		if i == 0 {
			return nl, nil
		}
		if nodes[i-1] == nil {
			return "null", nil
		}
		s, err := nodes[i-1].Semantic()
		if err != nil {
			return "", err
		}
		return s, nil
	default:
		return "", fmt.Errorf("arg.Type: %s invalid", arg.Type)
	}
}
