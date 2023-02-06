package main

import "fmt"

type Node interface {
	compile()
	String() string
}

type If struct {
	Cond Node
	Then []Node
	Else []Node
	Line int
}

func (i If) String() string {
	return fmt.Sprintf("If(%s, %s, %s)", i.Cond.String(), i.Then, i.Else)
}

type While struct {
	Cond Node
	Body []Node
	Line int
}

func (w While) String() string {
	return fmt.Sprintf("While(%s, %s)", w.Cond.String(), w.Body)
}

type For struct {
	Name string
	In   Node
	Body []Node
	Line int
}

func (f For) String() string {
	return fmt.Sprintf("For(%s, %s, %s)", f.Name, f.In.String(), f.Body)
}

type Func struct {
	Name string
	Args []string
	Body []Node
	Line int
}

func (f Func) String() string {
	return fmt.Sprintf("Func(%s, %s, %s)", f.Name, f.Args, f.Body)
}

type Sub struct {
	Name string
	Args []string
	Body []Node
	Line int
}

func (s Sub) String() string {
	return fmt.Sprintf("Sub(%s, %s, %s)", s.Name, s.Args, s.Body)
}

type Return struct {
	Expr Node
	Line int
}

func (r Return) String() string {
	return fmt.Sprintf("Return(%s)", r.Expr.String())
}

type Break struct {
	Line int
}

func (b Break) String() string {
	return "Break()"
}

type Continue struct {
	Line int
}

func (c Continue) String() string {
	return "Continue()"
}

type FuncCall struct {
	Func Node
	Args []Node
	Line int
}

func (f FuncCall) String() string {
	return fmt.Sprintf("FuncCall(%s, %s)", f.Func.String(), f.Args)
}

//go:generate stringer -type=OperatorType
type OperatorType int

const (
	Plus OperatorType = iota
	Minus
	Multiply
	Divide
	Modulo
	Equal
	NotEqual
	Less
	LessEqual
	Greater
	GreaterEqual
	And
	Or
	Not
	Assign
	PlusAssign
	MinusAssign
	MultiplyAssign
	DivideAssign
)

type Assignment struct {
	Left  string
	Right Node
	Op    OperatorType
	Line  int
}

func (a Assignment) String() string {
	return fmt.Sprintf("Assignment(%s, %s, %s)", a.Left, a.Op.String(), a.Right.String())
}

type AssignIndex struct {
	Left  string
	Index Node
	Right Node
	Op    OperatorType
	Line  int
}

func (a AssignIndex) String() string {
	return fmt.Sprintf("AssignIndex(%s, %s, %s, %s)", a.Left, a.Index.String(), a.Op.String(), a.Right.String())
}

type Unary struct {
	Expr Node
	Op   OperatorType
	Line int
}

func (u Unary) String() string {
	return fmt.Sprintf("Unary(%s, %s)", u.Op.String(), u.Expr.String())
}

type Binary struct {
	Left  Node
	Right Node
	Op    OperatorType
	Line  int
}

func (b Binary) String() string {
	return fmt.Sprintf("Binary(%s, %s, %s)", b.Left.String(), b.Op.String(), b.Right.String())
}

type Number struct {
	Value float64
	Line  int
}

func (n Number) String() string {
	return fmt.Sprintf("Number(%f)", n.Value)
}

type String struct {
	Value string
	Line  int
}

func (s String) String() string {
	return fmt.Sprintf("String(%q)", s.Value)
}

type Bool struct {
	Value bool
	Line  int
}

func (b Bool) String() string {
	return fmt.Sprintf("Bool(%t)", b.Value)
}

type List struct {
	Values []Node
	Line   int
}

func (l List) String() string {
	return fmt.Sprintf("List(%s)", l.Values)
}

type Variable struct {
	Name string
	Line int
}

func (v Variable) String() string {
	return fmt.Sprintf("Variable(%s)", v.Name)
}

type Index struct {
	Expr  Node
	Index Node
	Line  int
}

func (i Index) String() string {
	return fmt.Sprintf("Index(%s, %s)", i.Expr.String(), i.Index.String())
}
