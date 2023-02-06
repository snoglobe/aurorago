package main

import (
	"fmt"
	"strconv"
)

type Parser struct {
	lexer   *Lexer
	current Token
}

func NewParser(lexer *Lexer) *Parser {
	parser := &Parser{lexer, Token{}}
	parser.current = parser.lexer.Next()
	return parser
}

func (p *Parser) eat(typ TokenType) Token {
	if p.current.Type == typ {
		toRet := p.current
		p.current = p.lexer.Next()
		return toRet
	} else if typ == NewlineTok && p.current.Type == EofTok {
		return p.current
	} else {
		panic(fmt.Sprintf("Expected %s, got %s", typ, p.current.Type))
	}
}

func (p *Parser) peek(typ TokenType) bool {
	return p.current.Type == typ
}

func (p *Parser) peekNext() TokenType {
	return p.current.Type
}

func (p *Parser) ifStatement() Node {
	line := p.eat(IfTok).Line
	cond := p.expression()
	if p.peek(NewlineTok) {
		p.eat(NewlineTok)
		stmts := make([]Node, 0)
		elseStmts := make([]Node, 0)
		for !p.peek(EndTok) && !p.peek(ElseTok) {
			stmts = append(stmts, p.statement())
		}
		if p.peek(ElseTok) {
			p.eat(ElseTok)
			p.eat(NewlineTok)
			for !p.peek(EndTok) {
				elseStmts = append(elseStmts, p.statement())
			}
		}
		p.eat(EndTok)
		p.eat(NewlineTok)
		return If{cond, stmts, elseStmts, line}
	} else {
		stmt := p.statement()
		elseStmt := make([]Node, 0)
		if p.peek(ElseTok) {
			p.eat(ElseTok)
			elseStmt = append(elseStmt, p.statement())
		}
		return If{cond, []Node{stmt}, elseStmt, line}
	}
}

func (p *Parser) whileStatement() Node {
	line := p.eat(WhileTok).Line
	cond := p.expression()
	if p.peek(NewlineTok) {
		p.eat(NewlineTok)
		stmts := make([]Node, 0)
		for !p.peek(EndTok) {
			stmts = append(stmts, p.statement())
		}
		p.eat(EndTok)
		p.eat(NewlineTok)
		return While{cond, stmts, line}
	} else {
		stmt := p.statement()
		return While{cond, []Node{stmt}, line}
	}
}

func (p *Parser) forStatement() Node {
	line := p.eat(ForTok).Line
	name := p.eat(IdTok).Value
	p.eat(CommaTok)
	iter := p.expression()
	if p.peek(NewlineTok) {
		p.eat(NewlineTok)
		stmts := make([]Node, 0)
		for !p.peek(EndTok) {
			stmts = append(stmts, p.statement())
		}
		p.eat(EndTok)
		p.eat(NewlineTok)
		return For{name, iter, stmts, line}
	} else {
		stmt := p.statement()
		return For{name, iter, []Node{stmt}, line}
	}
}

func (p *Parser) functionStatement() Node {
	line := p.eat(FnTok).Line
	name := p.eat(IdTok).Value
	args := make([]string, 0)
	if !p.peek(NewlineTok) && !p.peek(ArrowTok) {
		for {
			args = append(args, p.eat(IdTok).Value)
			if p.peek(NewlineTok) || p.peek(ArrowTok) {
				break
			}
			p.eat(CommaTok)
		}
	}
	if p.peek(NewlineTok) {
		p.eat(NewlineTok)
		stmts := make([]Node, 0)
		for !p.peek(EndTok) {
			stmts = append(stmts, p.statement())
		}
		p.eat(EndTok)
		p.eat(NewlineTok)
		return Func{name, args, stmts, line}
	} else {
		p.eat(ArrowTok)
		expr := p.expression()
		p.eat(NewlineTok)
		return Func{name, args, []Node{Return{expr, line}}, line}
	}
}

func (p *Parser) subStatement() Node {
	line := p.eat(SubTok).Line
	name := p.eat(IdTok).Value
	args := make([]string, 0)
	if !p.peek(NewlineTok) && !p.peek(ArrowTok) {
		for {
			args = append(args, p.eat(IdTok).Value)
			if p.peek(NewlineTok) || p.peek(ArrowTok) {
				break
			}
			p.eat(CommaTok)
		}
	}
	if p.peek(NewlineTok) {
		p.eat(NewlineTok)
		stmts := make([]Node, 0)
		for !p.peek(EndTok) {
			stmts = append(stmts, p.statement())
		}
		p.eat(EndTok)
		p.eat(NewlineTok)
		return Sub{name, args, stmts, line}
	} else {
		p.eat(ArrowTok)
		stmt := p.statement()
		p.eat(NewlineTok)
		return Sub{name, args, []Node{stmt}, line}
	}
}

func (p *Parser) returnStatement() Node {
	line := p.eat(ReturnTok).Line
	expr := p.expression()
	p.eat(NewlineTok)
	return Return{expr, line}
}

func (p *Parser) breakStatement() Node {
	line := p.eat(BreakTok).Line
	p.eat(NewlineTok)
	return Break{line}
}

func (p *Parser) continueStatement() Node {
	line := p.eat(ContinueTok).Line
	p.eat(NewlineTok)
	return Continue{line}
}

func (p *Parser) statement() Node {
	switch p.peekNext() {
	case IfTok:
		return p.ifStatement()
	case WhileTok:
		return p.whileStatement()
	case ForTok:
		return p.forStatement()
	case FnTok:
		return p.functionStatement()
	case SubTok:
		return p.subStatement()
	case ReturnTok:
		return p.returnStatement()
	case BreakTok:
		return p.breakStatement()
	case ContinueTok:
		return p.continueStatement()
	case IdTok:
		name := p.eat(IdTok)
		// TODO more assignment types
		if p.peek(AssignTok) {
			p.eat(AssignTok)
			expr := p.expression()
			p.eat(NewlineTok)
			return Assignment{name.Value, expr, Assign, name.Line}
		} else if p.peek(ColonTok) {
			p.eat(ColonTok)
			index := p.expression()
			// TODO more assignment types
			p.eat(AssignTok)
			expr := p.expression()
			p.eat(NewlineTok)
			return AssignIndex{name.Value, index, expr, Assign, name.Line}
		} else {
			args := make([]Node, 0)
			if !p.peek(NewlineTok) && !p.peek(EofTok) {
				for {
					args = append(args, p.expression())
					if p.peek(NewlineTok) || p.peek(EofTok) {
						break
					}
					p.eat(CommaTok)
				}
			}
			p.eat(NewlineTok)
			return FuncCall{Variable{name.Value, name.Line}, args, name.Line}
		}
	case NewlineTok:
		p.eat(NewlineTok)
		return p.statement()
	default:
		panic("Unexpected token: " + p.peekNext().String())
	}
}

func (p *Parser) primary() Node {
	switch p.peekNext() {
	case NumberTok:
		numToken := p.eat(NumberTok)
		v, _ := strconv.ParseFloat(numToken.Value, 64)
		return Number{v, numToken.Line}
	case StringTok:
		strToken := p.eat(StringTok)
		return String{strToken.Value, strToken.Line}
	case IdTok:
		idToken := p.eat(IdTok)
		return Variable{idToken.Value, idToken.Line}
	case TrueTok:
		val := p.eat(TrueTok)
		return Bool{true, val.Line}
	case FalseTok:
		val := p.eat(FalseTok)
		return Bool{false, val.Line}
	case LparenTok:
		p.eat(LparenTok)
		expr := p.expression()
		p.eat(RparenTok)
		return expr
	case LbraceTok:
		line := p.eat(LbraceTok).Line
		exprs := make([]Node, 0)
		if !p.peek(RbraceTok) {
			for {
				exprs = append(exprs, p.expression())
				if p.peek(RbraceTok) {
					break
				}
				p.eat(CommaTok)
			}
		}
		p.eat(RbraceTok)
		return List{exprs, line}
	default:
		panic("Unexpected token: " + p.peekNext().String())
	}
}

func (p *Parser) call() Node {
	expr := p.primary()
	for {
		switch p.peekNext() {
		case LparenTok:
			line := p.eat(LparenTok).Line
			args := make([]Node, 0)
			if !p.peek(RparenTok) {
				for {
					args = append(args, p.expression())
					if p.peek(RparenTok) {
						break
					}
					p.eat(CommaTok)
				}
			}
			p.eat(RparenTok)
			expr = FuncCall{expr, args, line}
		case ColonTok:
			line := p.eat(ColonTok).Line
			index := p.expression()
			expr = Index{expr, index, line}
		default:
			return expr
		}
	}
}

func (p *Parser) unary() Node {
	switch p.peekNext() {
	case MinusTok:
		line := p.eat(MinusTok).Line
		return Unary{p.unary(), Minus, line}
	case NotTok:
		line := p.eat(NotTok).Line
		return Unary{p.unary(), Not, line}
	default:
		return p.call()
	}
}

func (p *Parser) factor() Node {
	expr := p.unary()
	for {
		switch p.peekNext() {
		case StarTok:
			line := p.eat(StarTok).Line
			expr = Binary{expr, p.unary(), Multiply, line}
		case SlashTok:
			line := p.eat(SlashTok).Line
			expr = Binary{expr, p.unary(), Divide, line}
		case PercentTok:
			line := p.eat(PercentTok).Line
			expr = Binary{expr, p.unary(), Modulo, line}
		default:
			return expr
		}
	}
}

func (p *Parser) term() Node {
	expr := p.factor()
	for {
		switch p.peekNext() {
		case PlusTok:
			line := p.eat(PlusTok).Line
			expr = Binary{expr, p.factor(), Plus, line}
		case MinusTok:
			line := p.eat(MinusTok).Line
			expr = Binary{expr, p.factor(), Minus, line}
		default:
			return expr
		}
	}
}

func (p *Parser) comparison() Node {
	expr := p.term()
	for {
		switch p.peekNext() {
		case GreaterTok:
			line := p.eat(GreaterTok).Line
			expr = Binary{expr, p.term(), Greater, line}
		case GreaterEqualTok:
			line := p.eat(GreaterEqualTok).Line
			expr = Binary{expr, p.term(), GreaterEqual, line}
		case LessTok:
			line := p.eat(LessTok).Line
			expr = Binary{expr, p.term(), Less, line}
		case LessEqualTok:
			line := p.eat(LessEqualTok).Line
			expr = Binary{expr, p.term(), LessEqual, line}
		default:
			return expr
		}
	}
}

func (p *Parser) equality() Node {
	expr := p.comparison()
	for {
		switch p.peekNext() {
		case EqualTok:
			line := p.eat(EqualTok).Line
			expr = Binary{expr, p.comparison(), Equal, line}
		case NotequalTok:
			line := p.eat(NotequalTok).Line
			expr = Binary{expr, p.comparison(), NotEqual, line}
		default:
			return expr
		}
	}
}

func (p *Parser) and() Node {
	expr := p.equality()
	for {
		switch p.peekNext() {
		case AndTok:
			line := p.eat(AndTok).Line
			expr = Binary{expr, p.equality(), And, line}
		default:
			return expr
		}
	}
}

func (p *Parser) or() Node {
	expr := p.and()
	for {
		switch p.peekNext() {
		case OrTok:
			line := p.eat(OrTok).Line
			expr = Binary{expr, p.and(), Or, line}
		default:
			return expr
		}
	}
}

func (p *Parser) expression() Node {
	return p.or()
}
