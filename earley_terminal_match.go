package fmr

import (
	"strings"

	"github.com/liuzl/ling"
)

func terminalMatch(term *Term, token *ling.Token) bool {
	if term == nil || token == nil || term.Type != Terminal {
		return false
	}
	t := gTokens.get(term.Value)
	if term.Meta == nil || t == nil {
		if term.Value == token.Text {
			return true
		}
	} else {
		flags, _ := term.Meta.(string)
		switch {
		case strings.Contains(flags, "l"):
			if t.Annotations[ling.Lemma] == token.Annotations[ling.Lemma] {
				return true
			}
		case strings.Contains(flags, "i"):
			if strings.ToLower(t.Annotations[ling.Norm]) ==
				strings.ToLower(token.Annotations[ling.Norm]) {
				return true
			}
		}
	}
	return false
}
