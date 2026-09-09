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
