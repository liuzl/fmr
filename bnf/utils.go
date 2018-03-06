package bnf

import (
	"unicode"
)

func isSpace(r rune) bool {
	return unicode.IsSpace(r)
}
func isWord(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}
func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}
