package fmr

import (
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
	if n.Value.Rb == nil || n.Value.Rb.F == nil {
		if n.p == nil {
			return nil
		}
		return map[string]interface{}{
			"type": n.Term().Value,
			"text": n.OriginalText(),
			"pos":  n.Pos(),
		}
	}

	subnodes := []interface{}{}
	for _, node := range n.Children {
		subnodes = append(subnodes, node.Tree())
	}
	return map[string]interface{}{
		"type":  n.Term().Value,
		"text":  n.OriginalText(),
		"pos":   n.Pos(),
		"nodes": subnodes,
	}
}
