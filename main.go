package main

import (
	"fmt"
	"regexp"
)

// THE LEXER : token converter of the prohram
type TokenType string

const (
	TokenKwVar   TokenType = "KW_VAR"   // var or this case = yaha rakh is how we declare variable
	TokenKwPrint TokenType = "KW_PRINT" // instead of print we use "dikhao"
	TokenId      TokenType = "ID"       // variable names
	TokenNumber  TokenType = "NUMBER"   // numbers
	TokenOp      TokenType = "OP"       // operators(+ - * /)
	TokenAssign  TokenType = "ASSIGN"   // assignment (=)
	TokenLParen  TokenType = "LPAREN"   // (
	TokenRParen  TokenType = "RPAREN"   // )
	TokenNewLine TokenType = "NEWLINE"  //  new line /n
	TokenEOF     TokenType = "EOF"      // end of line
)

// Token struct:
/*
  A token strutc is basically a type , value and line
  type: TokenType
  Value : Actual text like dikhao or yaha rakh
  Line: this token at which line number ( to track or so)
*/

type Token struct {
	Type  TokenType
	Value string
	Line  int
}

// TokenDefinition maps the regex rules to our TokenType so
// some regex which is mapped onto toketype
// example tokenKwVAr = 'yaha'
type TokenDef struct {
	Type  TokenType
	Regex *regexp.Regexp
}

// Tokeniser: starts from the first pointer
// looks at each input runs it against a loop of regex and stores as slice of tyope tokendef
func tokenise(source string) []Token {
	defs := []TokenDef{
		{TokenKwVar, regexp.MustCompile(`^yaha`)},
		{TokenKwPrint, regexp.MustCompile(`^dikhao`)},
		{TokenNumber, regexp.MustCompile(`^\d+(\.\d+)?`)},
		{TokenAssign, regexp.MustCompile(`^=`)},
		{TokenOp, regexp.MustCompile(`^[\+\-\*/]`)},
		{TokenLParen, regexp.MustCompile(`^\(`)},
		{TokenRParen, regexp.MustCompile(`^\)`)},
		{TokenId, regexp.MustCompile(`^[a-zA-Z_]\w*`)},
	}
	fmt.Print(defs)

}

func main() {
	fmt.Println("Welcome to the world of MoLang(.mo)")
}
