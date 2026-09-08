package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
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

// 	Token struct:
//  A token strutc is basically a type , value and line
//  type: TokenType
//  Value : Actual text like dikhao or yaha rakh
//  Line: this token at which line number ( to track or so)

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

// PARSER & AST Formation (recursive descent parsing)
/*    Evaluating yaha distance = speed * time
		results in
	          VarDeclNode (yaha)
           /         \
   VarName: distance   Expr: BinOpNode (+)
                             /         \
               BinOpNode (*)             NumNode (15)
               /          \
       VarNode(speed)   VarNode(time)

*/

// Base node for all nodes type
type Node interface{}

// Print keyword node
type PrintNode struct {
	Expr Node
}
type VarDecNode struct {
	VarName string
	Expr    Node
}

type BinOpNode struct {
	Left  Node
	Op    string
	Right Node
}

type StringNode struct {
	Value string
}

type NumNode struct {
	Value float64
}

type VarNode struct {
	Name string
}

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func (p *Parser) current() Token {
	return p.tokens[p.pos]
}

func (p *Parser) consume(tType TokenType) Token {
	tok := p.current()
	if tok.Type != tType {
		panic(fmt.Sprintf("Bhai galat token nhi chalega yaha,chahiye %s par mila %s yaha par %d", tType, tok.Type, tok.Line))
	}
	p.pos++
	return tok
}

// Parse function return the slice of valid tokens of type Node
func (p *Parser) Parse() []Node {
	var statements []Node

	for p.current().Type != TokenEOF {
		if p.current().Type == TokenNewLine {
			p.consume(TokenNewLine)
			continue
		}
		statements = append(statements, p.parseStaements())
	}
	return statements
}

func (p *Parser) parseStaements() Node {
	tok := p.current()

	switch tok.Type {
	case TokenKwVar:
		p.consume(TokenKwVar)
		varName := p.consume(TokenId).Value
		p.consume(TokenAssign)
		expr := p.parseExpression()
		return &VarDecNode{VarName: varName, Expr: expr}

	case TokenKwPrint:
		p.consume(TokenKwPrint)
		p.consume(TokenLParen)
		expr := p.parseExpression()
		p.consume(TokenRParen)
		return &PrintNode{Expr: expr}
	default:
		panic(fmt.Sprintf("bhaya ye mujhe na pata kya syntax h %s yaha pe %d", tok.Value, tok.Line))

	}

}

func (p *Parser) parseExpression() Node {
	left := p.parsePrimary()
	for p.current().Type == TokenOp {
		op := p.consume(TokenOp).Value
		right := p.parsePrimary()
		left = &BinOpNode{Left: left, Op: op, Right: right}
	}
	return left
}

func (p *Parser) parsePrimary() Node {
	tok := p.current()
	switch tok.Type {
	case TokenNumber:
		p.consume(TokenNumber)
		val, _ := strconv.ParseFloat(tok.Value, 64)
		return &NumNode{Value: val}
	case TokenString:
		p.consume(TokenString)
		// Strip the wrapping double quotes from the string value before saving
		cleanVal := strings.Trim(tok.Value, "\"")
		return &StringNode{Value: cleanVal}
	case TokenId:
		p.consume(TokenId)
		return &VarNode{Name: tok.Value}
	case TokenLParen:
		p.consume(TokenLParen)
		expr := p.parseExpression()
		p.consume(TokenRParen)
		return expr
	default:
		panic(fmt.Sprintf("Bhaya hame number, string, variable ya expression chahiye, mila '%s' yaha %d", tok.Value, tok.Line))
	}
}

// Runtime  Environment (interepreter logic)

// variables store the global memory of the program
// eg: lets take a statement `yaha speed = 10*23`
// variable["speed"] =value
type Environment struct {
	variables map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{variables: make(map[string]any)}
}

func (e *Environment) Execute(nodes []Node) {
	for _, node := range nodes {
		e.evaluate(node)
	}
}

func (e *Environment) evaluate(node Node) any {
	switch n := node.(type) {
	case *NumNode:
		return n.Value
	case *StringNode:
		return n.Value
	case *VarNode:
		val, exist := e.variables[n.Name]
		if !exist {
			panic(fmt.Sprintf("Runtime Error: Variable '%s' initialization dhoond nahi paye!", n.Name))
		}
		return val

	case *VarDecNode:
		val := e.evaluate(n.Expr)
		e.variables[n.VarName] = val
		fmt.Printf("Environment: %s = %v\n", n.VarName, val)
		return val

	case *PrintNode:
		val := e.evaluate(n.Expr)
		fmt.Println(val)
		return val
	case *BinOpNode:
		leftVal := e.evaluate(n.Left).(float64)
		rightVal := e.evaluate(n.Right).(float64)
		switch n.Op {
		case "+":
			return leftVal + rightVal
		case "-":
			return leftVal - rightVal
		case "*":
			return leftVal * rightVal
		case "/":
			if rightVal == 0 {
				panic("runtime error: bhai 0 se kaise divide karega satvi fail")
			}
			return leftVal / rightVal

		}
	default:
		panic(fmt.Sprintf("Unknown AST node: %T", node))
	}
	return 0
}

func main() {
	fmt.Println("Welcome to the world of MoLang(.mo)")
	// Read the file path and parse them and pass to tokenise function
	// because this could be step 1 of our language to be standalone

	// timestart := time.Now()

	// check if user has provided the filepath or not
	if len(os.Args) < 2 {
		fmt.Println("Galti: file to do kaam karne ke liye !")
		fmt.Println("Usage: go run main.go <filename.mo> !")
		os.Exit(1)
	}
	filePath := os.Args[1]
	// ensure the extension of our file is .mo
	if filepath.Ext(filePath) != ".mo" {
		fmt.Println("Galti: bhaya kewal .mo wale hi file supported h!")
		os.Exit(1)
	}

	// read the file directly
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Maafi: Bhaya file khol na paaye ham!")
		os.Exit(1)
	}

	// convert this file bytes back to string
	fileString := string(fileBytes)

	lexer := tokenise(fileString)

	parseit := NewParser(lexer)
	ast := parseit.Parse()

	for _, node := range ast {
		fmt.Printf("%#v\n", node)
	}
	// timss := time.Since(timestart)
	// fmt.Println(lexer)
	// fmt.Print("toal time", timss.Seconds())
	env := NewEnvironment()
	env.Execute(ast)

}
