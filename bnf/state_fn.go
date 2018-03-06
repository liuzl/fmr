package bnf

func lexItem(l *lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case r == '{':
			l.emit(itemLeftBrace)
		case r == '}':
			l.emit(itemRightBrace)
		case r == '(':
			l.emit(itemLeftParenthese)
		case r == ')':
			l.emit(itemRightParenthese)
		case r == '=':
			l.emit(itemEqualSign)
		case r == '|':
			l.emit(itemOr)
		case r == '*':
			l.emit(itemStar)
		case r == '+':
			l.emit(itemPlus)
		case r == '?':
			l.emit(itemOpt)
		case r == '"':
			return lexString
		case isSpace(r):
			l.ignore()
		case isWord(r):
			l.backup()
			return lexWord
		case r == eof:
			l.emit(itemEOF)
			break Loop
		default:
			l.emit(itemChar)
		}
	}
	return nil
}

func lexString(l *lexer) stateFn {
	var prev rune
	for {
		switch r := l.next(); {
		case r == '"' && prev != '\\':
			l.emit(itemString)
			return lexItem
		case r == eof:
			return l.errorf("unterminated string")
		case prev == '\\' && r == '\\':
			prev = 0
		default:
			prev = r
		}
	}
}

func lexWord(l *lexer) stateFn {
	first := true
	for {
		switch r := l.next(); {
		case isWord(r):
		case isDigit(r) && !first:
		default:
			l.backup()
			l.emit(itemIdentifier)
			return lexItem
		}
		first = false
	}
}
