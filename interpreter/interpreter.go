package interpreter

import (
	"fmt"

	"molang/ast"
)

type Environment struct {
	variables map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{variables: make(map[string]any)}
}

func (e *Environment) Execute(program ast.Program) {
	for _, node := range program {
		e.evaluate(node)
	}
}

func (e *Environment) evaluate(node ast.Node) any {
	switch current := node.(type) {
	case *ast.NumberNode:
		return current.Value
	case *ast.StringNode:
		return current.Value
	case *ast.VariableNode:
		value, exists := e.variables[current.Name]
		if !exists {
			panic(fmt.Sprintf("variable %q is not initialized", current.Name))
		}
		return value
	case *ast.VarDeclNode:
		value := e.evaluate(current.Expr)
		e.variables[current.VarName] = value
		return value
	case *ast.PrintNode:
		value := e.evaluate(current.Expr)
		fmt.Println(value)
		return value
	case *ast.BinaryOpNode:
		return e.evaluateBinary(current)
	case *ast.UnaryOpNode:
		return e.evaluateUnary(current)
	default:
		panic(fmt.Sprintf("unknown AST node: %T", node))
	}
}

func (e *Environment) evaluateBinary(node *ast.BinaryOpNode) float64 {
	left, leftOK := e.evaluate(node.Left).(float64)
	right, rightOK := e.evaluate(node.Right).(float64)
	if !leftOK || !rightOK {
		panic("arithmetic operators require number operands")
	}
	switch node.Op {
	case "+":
		return left + right
	case "-":
		return left - right
	case "*":
		return left * right
	case "/":
		if right == 0 {
			panic("cannot divide by zero")
		}
		return left / right
	default:
		panic(fmt.Sprintf("unknown operator %q", node.Op))
	}
}

func (e *Environment) evaluateUnary(node *ast.UnaryOpNode) float64 {
	value, ok := e.evaluate(node.Operand).(float64)
	if !ok {
		panic("unary op requier a number operand")
	}

	switch node.Op {
	case "+":
		return value
	case "-":
		return -value
	default:
		panic(fmt.Sprintf("unknown value h bhai %q", node.Op))
	}
}
