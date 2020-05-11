package fmr

import (
	"fmt"
	"math/big"
)

// Eval returns the denotation of Node n
func (n *Node) Eval() (interface{}, error) {
	if n.Value.Rb == nil || n.Value.Rb.F == nil {
		if n.p == nil {
			return "", nil
		}
		return n.OriginalText(), nil
	}
	return n.fmrEval(n.Value.Rb.F, n.Children)
}

func (n *Node) fmrEval(f *FMR, children []*Node) (interface{}, error) {
	if f == nil {
		return "", nil
	}
	if f.Fn == "nf.I" {
		if len(f.Args) != 1 {
			return "", fmt.Errorf("the length of Args of nf.I should be one")
		}
		s, err := n.semEval(f.Args[0], children)
		if err != nil {
			return "", err
		}
		return s, nil
	}

	var args []interface{}
	for _, arg := range f.Args {
		s, err := n.semEval(arg, children)
		if err != nil {
			return "", err
		}
		args = append(args, s)
	}
	if Debug {
		fmt.Printf("funcs.Call(%s, %+v)\n", f.Fn, args)
	}
	return Call(f.Fn, args...)
}

func (n *Node) semEval(arg *Arg, nodes []*Node) (interface{}, error) {
	if arg == nil {
		return "", fmt.Errorf("arg is nil")
	}
	switch arg.Type {
	case "string":
		if s, ok := arg.Value.(string); ok {
			return s, nil
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
			return n.fmrEval(fmr, nodes)
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
			return n.NL(), nil
		}
		s, err := nodes[i-1].Eval()
		if err != nil {
			return "", err
		}
		return s, nil
	case "context":
		subnodes := []map[string]interface{}{}
		for _, node := range nodes {
			ni, err := node.Eval()
			if err != nil {
				ni = node.OriginalText()
			}
			subnodes = append(subnodes, map[string]interface{}{node.Term().Value: ni})
		}
		ret := map[string]interface{}{
			"text":  n.OriginalText(),
			"pos":   n.Pos(),
			"nodes": subnodes,
		}
		if n.Term().Type != Terminal {
			ret["type"] = n.Term().Value
		}
		return ret, nil
	default:
		return "", fmt.Errorf("arg.Type: %s invalid", arg.Type)
	}
}
