package bnf

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/liuzl/goutil"
)

// Semantic returns the stringified FMR of Node n
func (n *Node) Semantic() (string, error) {
	if n.Value.Rb == nil || n.Value.Rb.F == nil {
		if n.p == nil {
			return "", nil
		}
		var s []string
		for i := n.Value.Start + 1; i <= n.Value.End; i++ {
			s = append(s, n.p.columns[i].token)
		}
		start := n.p.columns[n.Value.Start+1].startByte
		end := n.p.columns[n.Value.End].endByte
		fmt.Printf("%s[%d:%d]:%s\n", n.p.text, start, end, n.p.text[start:end])
		return strconv.Quote(goutil.Join(s)), nil
	}
	return fmrStr(n.Value.Rb.F, n.Children)
}

func fmrStr(f *FMR, children []*Node) (string, error) {
	if f == nil {
		return "", nil
	}
	if f.Fn == "nf.I" {
		if len(f.Args) != 1 {
			return "", fmt.Errorf("the length of Args of nf.I should be one")
		}
		s, err := semStr(f.Args[0], children)
		if err != nil {
			return "", err
		}
		return s, nil
	}

	var args []string
	for _, arg := range f.Args {
		s, err := semStr(arg, children)
		if err != nil {
			return "", err
		}
		args = append(args, s)
	}
	return fmt.Sprintf("%s(%s)", f.Fn, strings.Join(args, ", ")), nil
}

func semStr(arg *Arg, nodes []*Node) (string, error) {
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
			return fmrStr(fmr, nodes)
		}
		return "", fmt.Errorf("arg.Value: %+v is not func", arg.Value)
	case "index":
		i, ok := arg.Value.(int)
		if !ok {
			return "", fmt.Errorf("arg.Value: %+v is not index", arg.Value)
		}
		if i < 1 || i > len(nodes) {
			return "", fmt.Errorf("i=%d not in range [1, %d]", i, len(nodes))
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
