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
	TokenString  TokenType = "STRING"   // strings
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
		{TokenString, regexp.MustCompile(`^"[^"\n]*"`)},
		{TokenAssign, regexp.MustCompile(`^=`)},
		{TokenOp, regexp.MustCompile(`^[\+\-\*/]`)},
		{TokenLParen, regexp.MustCompile(`^\(`)},
		{TokenRParen, regexp.MustCompile(`^\)`)},
		{TokenId, regexp.MustCompile(`^[a-zA-Z_]\w*`)},
	}

	var tokens []Token
	line := 1

	i := 0
	for i < len(source) {
		matchFound := false

		// 1. check for NewLine char '\n'
		if source[i] == '\n' {
			tokens = append(tokens, Token{Type: TokenNewLine, Value: "\n", Line: line})
			line++
			i++
			continue
		}

		// 2. skip white space in between (structural)
		if source[i] == ' ' || source[i] == '\t' || source[i] == '\r' {
			i++
			continue
		}

		// 3. Now start checking the actual string chars
		remaining := source[i:]
		for _, def := range defs {
			loc := def.Regex.FindStringIndex(remaining)
			if loc != nil && loc[0] == 0 { // Match found at start only
				matchStr := remaining[loc[0]:loc[1]]                                        // the actual matched string ( as defined in the TokenType)
				tokens = append(tokens, Token{Type: def.Type, Value: matchStr, Line: line}) // append to tokens of type Token
				i += len(matchStr)                                                          // move i to the next location after this match eg: yaha x = 10 so i now point at x
				matchFound = true
				break
			}
		}
		if !matchFound {
			panic(fmt.Sprintf("Mazak udao ge mere lang ka invalid ops kar ke!: %c at line %d", source[i], line))
		}
	}

	tokens = append(tokens, Token{Type: TokenEOF, Value: "", Line: line})

	return tokens

}

func main() {
	fmt.Println("Welcome to the world of MoLang(.mo)")
}
