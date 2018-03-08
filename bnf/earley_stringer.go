package bnf

import (
	"fmt"
	"os"
)

func (self *TableState) String() string {
	s := ""
	for i, term := range self.rb.Terms {
		if i == self.dot {
			s += DOT + " "
		}
		s += term.Value + " "
	}
	if self.dot == len(self.rb.Terms) {
		s += DOT
	}
	return fmt.Sprintf("%-6s -> %-20s [%d-%d]",
		self.name, s, self.start.index, self.end.index)
}

func (self *TableColumn) String() string {
	out := fmt.Sprintf("[%d] '%s'\n", self.index, self.token)
	out += "=======================================\n"
	for _, s := range self.states {
		out += s.String() + "\n"
	}
	return out
}

func (self *Parser) String() string {
	out := ""
	for _, c := range self.columns {
		out += c.String() + "\n"
	}
	return out
}

func (self *Node) Print(out *os.File) {
	self.PrintLevel(out, 0)
}

func (self *Node) PrintLevel(out *os.File, level int) {
	indentation := ""
	for i := 0; i < level; i++ {
		indentation += "  "
	}
	fmt.Fprintf(out, "%s%v\n", indentation, self.value)
	for _, child := range self.children {
		child.PrintLevel(out, level+1)
	}
}

func (self *Node) String() string {
	if len(self.children) > 0 {
		return fmt.Sprintf("%+v %+v", self.value, self.children)
	} else {
		return fmt.Sprintf("%+v", self.value)
	}
}
