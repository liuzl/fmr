package fmr

import (
	"fmt"

	"github.com/xlab/treeprint"
)

// TreePrint to out
func (n *Node) TreePrint() {
	tree := treeprint.New()
	tree.SetValue(n.Value)
	for _, child := range n.Children {
		tree.AddNode(child.Value)
	}
	fmt.Println(tree.String())
}
