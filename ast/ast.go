package ast

type Node interface{}

type Program []Node

type PrintNode struct {
	Expr Node
}

type VarDeclNode struct {
	VarName string
	Expr    Node
}

type BinaryOpNode struct {
	Left  Node
	Op    string
	Right Node
}

type StringNode struct {
	Value string
}

type NumberNode struct {
	Value float64
}

type VariableNode struct {
	Name string
}

// handles unary ops (dikhao(-2*10)-> -10)
type UnaryOpNode struct {
	Op      string
	Operand Node
}
