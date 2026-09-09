package main

import (
	"fmt"
	"os"
	"path/filepath"

	"molang/interpreter"
	"molang/lexer"
	"molang/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run . <filename.mo>")
		os.Exit(1)
	}
	if filepath.Ext(os.Args[1]) != ".mo" {
		fmt.Println("only .mo files are supported")
		os.Exit(1)
	}

	source, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("could not read %s: %v\n", os.Args[1], err)
		os.Exit(1)
	}

	program := parser.New(lexer.Tokenize(string(source))).Parse()
	interpreter.NewEnvironment().Execute(program)
}
