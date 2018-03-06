package bnf

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

type item struct {
	typ itemType
	pos int
	val string
}

func (i *item) String() string {
	u, _ := strconv.Unquote(i.val)
	return fmt.Sprintf("%s,%d,%s,%s", i.typ, i.pos, i.val, u)
}

type itemType int

const (
	itemError itemType = iota
	itemEOF
	itemIdentifier
	itemLeftBrace       // '{'
	itemRightBrace      // '}'
	itemLeftParenthese  // '('
	itemRightParenthese // ')'
	itemEqualSign       // '='
	itemOr              // '|'
	itemStar            // '*'
	itemPlus            // '+'
	itemOpt             // '?'
	itemString
	itemChar
)

const eof = -1

type stateFn func(*lexer) stateFn

type lexer struct {
	input string
	start int
	pos   int
	width int
	state stateFn
	items chan item
}

func (l *lexer) next() rune {
	if l.isEOF() {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += w
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) isEOF() bool {
	return l.pos >= len(l.input)
}

func (l *lexer) contains(c string) bool {
	return strings.Index(l.input[l.pos:], c) >= 0
}

func (l *lexer) acceptUntil(c string) {
	for !strings.ContainsRune(c, l.next()) {
	}
	l.backup()
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.pos, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) run() {
	for l.state = lexItem; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item),
	}

	go l.run()
	return l
}
