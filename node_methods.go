package fmr

import (
	"fmt"
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
	if n.Value == nil { //|| n.Value.Rb == nil || len(n.Value.Rb.Terms) < 1 {
		return nil
	}
	if n.Value.Term.Value == GammaRule {
		return n.Value.Rb.Terms[0]
	}
	return n.Value.Term
}

// F returns the FMR signature of node
func (n *Node) F() *FMR {
	if n.Value == nil || n.Value.Rb == nil || len(n.Value.Rb.Terms) < 1 {
		return nil
	}
	if n.Value.Term.Value == GammaRule {
		return n.Children[0].Value.Rb.F
	}
	return n.Value.Rb.F
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

// Tree returns the parsed tree of Node n
func (n *Node) Tree() map[string]interface{} {
	if n.Value.Term.Value == GammaRule {
		return n.Children[0].Tree()
	}
	if n.p == nil {
		return nil
	}
	ret := map[string]interface{}{
		"type": n.Term().Value,
		"text": n.OriginalText(),
		"pos":  n.Pos(),
		//"fmr":  n.F(),
	}
	if n.Value.Rb == nil || n.Value.Rb.F == nil {
		return ret
	}
	if len(n.Children) == 1 {
		s := n.Children[0].Value.Term.Value
		if strings.HasPrefix(s, "g_t_") || strings.HasPrefix(s, "l_t_") {
			return ret
		}
	}

	subnodes := []interface{}{}
	for _, node := range n.Children {
		subnodes = append(subnodes, node.Tree())
	}
	ret["nodes"] = subnodes
	return ret
}

// Bracketed returns the labeled bracket notation of Node
func (n *Node) Bracketed() string {
	if n.Value.Term.Value == GammaRule {
		return n.Children[0].Bracketed()
	}
	if n.p == nil {
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, "[%s ", n.Term().Value)
	if n.Value.Rb == nil || n.Value.Rb.F == nil {
		fmt.Fprintf(&b, "%s]", strconv.Quote(n.OriginalText()))
		return b.String()
	}
	if len(n.Children) == 1 {
		s := n.Children[0].Value.Term.Value
		if strings.HasPrefix(s, "g_t_") || strings.HasPrefix(s, "l_t_") {
			fmt.Fprintf(&b, "%s]", strconv.Quote(n.OriginalText()))
			return b.String()
		}
	}
	subnodes := []string{}
	for _, node := range n.Children {
		subnodes = append(subnodes, node.Bracketed())
	}
	fmt.Fprintf(&b, "%s]", strings.Join(subnodes, " "))
	return b.String()
}
