package main

import "fmt"

type TokenType int

//go:generate stringer -type=TokenType
const (
	IdTok TokenType = iota
	NumberTok
	StringTok

	NewlineTok

	LparenTok
	RparenTok
	LbraceTok
	RbraceTok
	ArrowTok
	CommaTok
	ColonTok

	PlusTok
	MinusTok
	StarTok
	SlashTok
	PercentTok
	EqualTok
	NotequalTok
	LessTok
	LessEqualTok
	GreaterTok
	GreaterEqualTok
	AssignTok
	PlusAssignTok
	MinusAssignTok
	StarAssignTok
	SlashAssignTok
	PercentAssignTok

	IfTok
	ElseTok
	WhileTok
	ForTok
	FnTok
	SubTok
	ReturnTok
	BreakTok
	ContinueTok
	TrueTok
	FalseTok
	AndTok
	OrTok
	NotTok
	EndTok

	EofTok
)

var keywords = map[string]TokenType{
	"if":       IfTok,
	"else":     ElseTok,
	"while":    WhileTok,
	"for":      ForTok,
	"fn":       FnTok,
	"sub":      SubTok,
	"return":   ReturnTok,
	"break":    BreakTok,
	"continue": ContinueTok,
	"true":     TrueTok,
	"false":    FalseTok,
	"and":      AndTok,
	"or":       OrTok,
	"not":      NotTok,
	"end":      EndTok,
}

type Token struct {
	Type  TokenType
	Value string
	Line  int
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%s, '%s', %d)", t.Type.String(), t.Value, t.Line)
}

type Lexer struct {
	input string
	start int
	pos   int
	line  int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input, 0, 0, 1}
}

func (l *Lexer) scanString() string {
	l.pos++
	start := l.pos
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '"' {
			l.pos++
			return l.input[start : l.pos-1]
		}
		l.pos++
	}
	panic(fmt.Sprintf("Unterminated string at line %d", l.line))
}

func (l *Lexer) scanNumber() string {
	start := l.pos
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch < '0' || ch > '9' {
			break
		}
		l.pos++
	}
	if l.pos < len(l.input) && l.input[l.pos] == '.' {
		l.pos++
		for l.pos < len(l.input) {
			ch := l.input[l.pos]
			if ch < '0' || ch > '9' {
				break
			}
			l.pos++
		}
	}
	return l.input[start:l.pos]
}

func (l *Lexer) scanId() string {
	start := l.pos
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if !isAlpha(ch) && !isDigit(ch) {
			break
		}
		l.pos++
	}
	return l.input[start:l.pos]
}

func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) Next() Token {
	hadError := false
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if hadError {
			panic(fmt.Sprintf("Unexpected character %c at line %d", ch, l.line))
		}
		switch ch {
		case ' ', '\t':
			l.pos++
		case '\n':
			l.pos++
			l.line++
			return Token{NewlineTok, "\n", l.line - 1}
		case '(':
			l.pos++
			return Token{LparenTok, "(", l.line}
		case ')':
			l.pos++
			return Token{RparenTok, ")", l.line}
		case '{':
			l.pos++
			return Token{LbraceTok, "{", l.line}
		case '}':
			l.pos++
			return Token{RbraceTok, "}", l.line}
		case ',':
			l.pos++
			return Token{CommaTok, ",", l.line}
		case ':':
			l.pos++
			return Token{ColonTok, ":", l.line}
		case '+':
			l.pos++
			if l.pos < len(l.input) && l.input[l.pos] == '=' {
				l.pos++
				return Token{PlusAssignTok, "+=", l.line}
			}
			return Token{PlusTok, "+", l.line}
		case '-':
			l.pos++
			if l.pos < len(l.input) && l.input[l.pos] == '=' {
				l.pos++
				return Token{MinusAssignTok, "-=", l.line}
			} else if l.pos < len(l.input) && l.input[l.pos] == '>' {
				l.pos++
				return Token{ArrowTok, "->", l.line}
			}
			return Token{MinusTok, "-", l.line}
		case '*':
			l.pos++
			if l.pos < len(l.input) && l.input[l.pos] == '=' {
				l.pos++
				return Token{StarAssignTok, "*=", l.line}
			}
			return Token{StarTok, "*", l.line}
		case '/':
			l.pos++
			if l.pos < len(l.input) && l.input[l.pos] == '=' {
				l.pos++
				return Token{SlashAssignTok, "/=", l.line}
			}
			return Token{SlashTok, "/", l.line}
		case '%':
			l.pos++
			if l.pos < len(l.input) && l.input[l.pos] == '=' {
				l.pos++
				return Token{PercentAssignTok, "%=", l.line}
			}
			return Token{PercentTok, "%", l.line}
		case '=':
			l.pos++
			if l.pos < len(l.input) && l.input[l.pos] == '=' {
				l.pos++
				return Token{EqualTok, "==", l.line}
			}
			return Token{AssignTok, "=", l.line}
		case '!':
			l.pos++
			if l.pos < len(l.input) && l.input[l.pos] == '=' {
				l.pos++
				return Token{NotequalTok, "!=", l.line}
			}
			hadError = true
		case '<':
			l.pos++
			if l.pos < len(l.input) && l.input[l.pos] == '=' {
				l.pos++
				return Token{LessEqualTok, "<=", l.line}
			}
			return Token{LessTok, "<", l.line}
		case '>':
			l.pos++
			if l.pos < len(l.input) && l.input[l.pos] == '=' {
				l.pos++
				return Token{GreaterEqualTok, ">=", l.line}
			}
			return Token{GreaterTok, ">", l.line}
		case '&':
			l.pos++
			if l.pos < len(l.input) && l.input[l.pos] == '&' {
				l.pos++
				return Token{AndTok, "&&", l.line}
			}
			hadError = true
		case '"':
			return Token{
				StringTok,
				l.scanString(),
				l.line,
			}
		default:
			if isDigit(ch) {
				return Token{
					NumberTok,
					l.scanNumber(),
					l.line,
				}
			}
			if isAlpha(ch) {
				ident := l.scanId()
				if tok, ok := keywords[ident]; ok {
					return Token{tok, ident, l.line}
				} else {
					return Token{IdTok, ident, l.line}
				}
			}
			hadError = true
		}
	}
	return Token{EofTok, "", l.line}
}
