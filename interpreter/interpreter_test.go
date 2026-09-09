package interpreter

import (
	"testing"

	"molang/ast"
)

func TestEvaluateArithmetic(t *testing.T) {
	environment := NewEnvironment()
	expression := &ast.BinaryOpNode{
		Left:  &ast.NumberNode{Value: 2},
		Op:    "+",
		Right: &ast.NumberNode{Value: 3},
	}

	result, ok := environment.evaluate(expression).(float64)
	if !ok {
		t.Fatalf("expected numeric result, got %T", environment.evaluate(expression))
	}
	if result != 5 {
		t.Fatalf("expected 5, got %v", result)
	}
}

func TestEvaluateUnaryOperators(t *testing.T) {
	environment := NewEnvironment()

	negative := environment.evaluate(&ast.UnaryOpNode{
		Op:      "-",
		Operand: &ast.NumberNode{Value: 5},
	})
	if negative != float64(-5) {
		t.Fatalf("expected -5, got %v", negative)
	}

	positive := environment.evaluate(&ast.UnaryOpNode{
		Op:      "+",
		Operand: &ast.NumberNode{Value: 5},
	})
	if positive != float64(5) {
		t.Fatalf("expected 5, got %v", positive)
	}
}

func TestEvaluateUnaryWithArithmetic(t *testing.T) {
	environment := NewEnvironment()
	expression := &ast.BinaryOpNode{
		Left: &ast.NumberNode{Value: 2},
		Op:   "*",
		Right: &ast.UnaryOpNode{
			Op:      "-",
			Operand: &ast.NumberNode{Value: 3},
		},
	}

	result := environment.evaluate(expression)
	if result != float64(-6) {
		t.Fatalf("expected -6, got %v", result)
	}
}

func TestEvaluateUnknownVariablePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected unknown variable to panic")
		}
	}()

	NewEnvironment().evaluate(&ast.VariableNode{Name: "missing"})
}
