package bnf

import (
	"fmt"
)

func Parse(d string) (*Grammar, error) {
	if d == "" {
		return nil, fmt.Errorf("parse error: nothing to parse")
	}
	l := lex(d)
	t := <-l.items
	if t.typ != itemIdentifier {
		return nil, fmt.Errorf("parse error: must begin with a grammar name")
	}
	g := &Grammar{Name: t.val}
	t = <-l.items
	if t.typ != itemLeftBrace {
		return nil, fmt.Errorf("parse error: '{' expected")
	}
	var rule *Rule
	var ident string
Loop:
	for {
		t = <-l.items
		switch t.typ {
		case itemRightBrace:
			break
		case itemEOF:
			break Loop
		case itemError:
			return nil, fmt.Errorf("parse error: %s", t.val)
		case itemIdentifier:
			if ident == "" {
				ident = t.val
			} else {
			}
		case itemEqualSign:
			if rule != nil {
				// save the pre rule
			}
		}
	}

	return g, nil
}
