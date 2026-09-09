package lexer

import (
	"fmt"
	"regexp"
)

type TokenType string

const (
	TokenKwVar   TokenType = "KW_VAR"
	TokenKwPrint TokenType = "KW_PRINT"
	TokenKwIf    TokenType = "KW_IF"
	TokenKwElse  TokenType = "KW_ELSE"
	TokenId      TokenType = "ID"
	TokenNumber  TokenType = "NUMBER"
	TokenString  TokenType = "STRING"
	TokenBoolean TokenType = "BOOLEAN"
	TokenCompare TokenType = "COMPARE"
	TokenOp      TokenType = "OP"
	TokenAssign  TokenType = "ASSIGN"
	TokenLParen  TokenType = "LPAREN"
	TokenRParen  TokenType = "RPAREN"
	TokenLBrace  TokenType = "LBRACE"
	TokenRBrace  TokenType = "RBRACE"
	TokenNewLine TokenType = "NEWLINE"

	TokenEOF TokenType = "EOF"
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
}

type tokenDefinition struct {
	typeName TokenType
	pattern  *regexp.Regexp
}

func Tokenize(source string) []Token {
	definitions := []tokenDefinition{
		{TokenKwIf, regexp.MustCompile(`^agar\b`)},
		{TokenKwElse, regexp.MustCompile(`^warna\b`)},
		{TokenBoolean, regexp.MustCompile(`^(true|false)\b`)},
		{TokenKwVar, regexp.MustCompile(`^yaha\b`)},
		{TokenKwPrint, regexp.MustCompile(`^dikhao\b`)},
		{TokenNumber, regexp.MustCompile(`^\d+(\.\d+)?`)},
		{TokenString, regexp.MustCompile(`^"[^"\n]*"`)},

		{TokenCompare, regexp.MustCompile(`^(==|!=|<=|>=|<|>)`)},
		{TokenAssign, regexp.MustCompile(`^=`)},

		{TokenOp, regexp.MustCompile(`^[+\-*/]`)},
		{TokenLParen, regexp.MustCompile(`^\(`)},
		{TokenRParen, regexp.MustCompile(`^\)`)},
		{TokenLBrace, regexp.MustCompile(`^\{`)},
		{TokenRBrace, regexp.MustCompile(`^\}`)},
		{TokenId, regexp.MustCompile(`^[a-zA-Z_]\w*`)},
	}

	var tokens []Token
	line := 1
	for position := 0; position < len(source); {
		if source[position] == '\n' {
			tokens = append(tokens, Token{Type: TokenNewLine, Value: "\n", Line: line})
			line++
			position++
			continue
		}
		if source[position] == ' ' || source[position] == '\t' || source[position] == '\r' {
			position++
			continue
		}

		remaining := source[position:]
		matched := false
		for _, definition := range definitions {
			match := definition.pattern.FindString(remaining)
			if match == "" {
				continue
			}
			tokens = append(tokens, Token{Type: definition.typeName, Value: match, Line: line})
			position += len(match)
			matched = true
			break
		}
		if !matched {
			panic(fmt.Sprintf("invalid character %q at line %d", source[position], line))
		}
	}

	return append(tokens, Token{Type: TokenEOF, Line: line})
}
