package parser

import (
	"testing"

	"molang/ast"
	"molang/lexer"
)

func TestParseProgram(t *testing.T) {
	program := New(lexer.Tokenize("yaha x = 10\ndikhao(x)")).Parse()

	if len(program) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program))
	}
	declaration, ok := program[0].(*ast.VarDeclNode)
	if !ok {
		t.Fatalf("expected variable declaration, got %T", program[0])
	}
	if declaration.VarName != "x" {
		t.Fatalf("expected variable name x, got %q", declaration.VarName)
	}
	if _, ok := declaration.Expr.(*ast.NumberNode); !ok {
		t.Fatalf("expected numeric initializer, got %T", declaration.Expr)
	}
	if _, ok := program[1].(*ast.PrintNode); !ok {
		t.Fatalf("expected print statement, got %T", program[1])
	}
}

func TestParseOperatorPrecedence(t *testing.T) {
	program := New(lexer.Tokenize("dikhao(2 + 3 * 4)")).Parse()
	printNode := program[0].(*ast.PrintNode)
	addition := printNode.Expr.(*ast.BinaryOpNode)

	if addition.Op != "+" {
		t.Fatalf("expected addition at the root, got %q", addition.Op)
	}
	multiplication, ok := addition.Right.(*ast.BinaryOpNode)
	if !ok || multiplication.Op != "*" {
		t.Fatalf("expected multiplication on the right side, got %#v", addition.Right)
	}
}

func TestParseUnaryExpression(t *testing.T) {
	program := New(lexer.Tokenize("dikhao(-5)")).Parse()
	printNode := program[0].(*ast.PrintNode)
	unary, ok := printNode.Expr.(*ast.UnaryOpNode)
	if !ok {
		t.Fatalf("expected unary expression, got %T", printNode.Expr)
	}
	if unary.Op != "-" {
		t.Fatalf("expected unary minus, got %q", unary.Op)
	}
	if number, ok := unary.Operand.(*ast.NumberNode); !ok || number.Value != 5 {
		t.Fatalf("expected unary operand to be number 5, got %#v", unary.Operand)
	}
}

func TestParseUnaryExpressionAfterMultiplication(t *testing.T) {
	program := New(lexer.Tokenize("dikhao(2 * -3)")).Parse()
	printNode := program[0].(*ast.PrintNode)
	multiplication, ok := printNode.Expr.(*ast.BinaryOpNode)
	if !ok || multiplication.Op != "*" {
		t.Fatalf("expected multiplication root, got %#v", printNode.Expr)
	}
	if unary, ok := multiplication.Right.(*ast.UnaryOpNode); !ok || unary.Op != "-" {
		t.Fatalf("expected unary minus on multiplication right side, got %#v", multiplication.Right)
	}
}

func TestParseNestedUnaryExpression(t *testing.T) {
	program := New(lexer.Tokenize("dikhao(--5)")).Parse()
	printNode := program[0].(*ast.PrintNode)
	outer, ok := printNode.Expr.(*ast.UnaryOpNode)
	if !ok || outer.Op != "-" {
		t.Fatalf("expected outer unary minus, got %#v", printNode.Expr)
	}
	inner, ok := outer.Operand.(*ast.UnaryOpNode)
	if !ok || inner.Op != "-" {
		t.Fatalf("expected nested unary minus, got %#v", outer.Operand)
	}
}
