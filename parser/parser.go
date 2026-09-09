package parser

import (
	"fmt"
	"strconv"
	"strings"

	"molang/ast"
	"molang/lexer"
)

type Parser struct {
	tokens   []lexer.Token
	position int
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() ast.Program {
	var statements ast.Program
	for p.current().Type != lexer.TokenEOF {
		if p.current().Type == lexer.TokenNewLine {
			p.consume(lexer.TokenNewLine)
			continue
		}
		statements = append(statements, p.parseStatement())
	}
	return statements
}

func (p *Parser) current() lexer.Token {
	return p.tokens[p.position]
}

func (p *Parser) consume(expected lexer.TokenType) lexer.Token {
	token := p.current()
	if token.Type != expected {
		panic(fmt.Sprintf("expected %s but found %s at line %d", expected, token.Type, token.Line))
	}
	p.position++
	return token
}

func (p *Parser) parseStatement() ast.Node {
	switch p.current().Type {
	case lexer.TokenKwVar:
		p.consume(lexer.TokenKwVar)
		name := p.consume(lexer.TokenId).Value
		p.consume(lexer.TokenAssign)
		return &ast.VarDeclNode{VarName: name, Expr: p.parseExpression()}
	case lexer.TokenKwPrint:
		p.consume(lexer.TokenKwPrint)
		p.consume(lexer.TokenLParen)
		expression := p.parseExpression()
		p.consume(lexer.TokenRParen)
		return &ast.PrintNode{Expr: expression}
	default:
		token := p.current()
		panic(fmt.Sprintf("unexpected token %q at line %d", token.Value, token.Line))
	}
}

func (p *Parser) parseExpression() ast.Node {
	return p.parseTerm()
}

func (p *Parser) parseTerm() ast.Node {
	left := p.parseFactor()
	for p.current().Type == lexer.TokenOp && (p.current().Value == "+" || p.current().Value == "-") {
		op := p.consume(lexer.TokenOp).Value
		right := p.parseFactor()
		left = &ast.BinaryOpNode{Left: left, Op: op, Right: right}
	}
	return left
}

func (p *Parser) parseFactor() ast.Node {
	left := p.parseUnary()
	for p.current().Type == lexer.TokenOp && (p.current().Value == "*" || p.current().Value == "/") {
		op := p.consume(lexer.TokenOp).Value
		right := p.parseUnary()
		left = &ast.BinaryOpNode{Left: left, Op: op, Right: right}
	}
	return left
}

func (p *Parser) parseUnary() ast.Node {
	if p.current().Type == lexer.TokenOp && (p.current().Value == "-" || p.current().Value == "+") {
		op := p.consume(lexer.TokenOp).Value

		return &ast.UnaryOpNode{
			Op:      op,
			Operand: p.parseUnary(),
		}
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() ast.Node {
	token := p.current()
	switch token.Type {
	case lexer.TokenNumber:
		p.consume(lexer.TokenNumber)
		value, _ := strconv.ParseFloat(token.Value, 64)
		return &ast.NumberNode{Value: value}
	case lexer.TokenString:
		p.consume(lexer.TokenString)
		return &ast.StringNode{Value: strings.Trim(token.Value, `"`)}
	case lexer.TokenId:
		p.consume(lexer.TokenId)
		return &ast.VariableNode{Name: token.Value}
	case lexer.TokenLParen:
		p.consume(lexer.TokenLParen)
		expression := p.parseExpression()
		p.consume(lexer.TokenRParen)
		return expression
	default:
		panic(fmt.Sprintf("expected an expression at line %d", token.Line))
	}
}
