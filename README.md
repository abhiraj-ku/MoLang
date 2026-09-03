# MoLang (.mo)

MoLang is a small toy programming language written from scratch in Go. Its
syntax is inspired by simple Hindi words:

- `yaha` declares a variable
- `dikhao` prints a value
- `+`, `-`, `*`, and `/` perform basic calculations

The project is being built step by step with a custom lexer, parser, and
execution environment. No external language framework is required.

## Example

The intended MoLang syntax looks like this:

```molang
yaha x = 10
yaha y = 5
dikhao x + y
```

This program declares two variables and prints their sum.

## Language Basics

### Variables

Use `yaha` followed by a variable name, `=`, and a number:

```molang
yaha age =  twenty
```

Variables can use letters, numbers, and underscores, but must start with a
letter or underscore.

### Printing

Use `dikhao` to print a number or calculation:

```molang
dikhao 10
dikhao 10 * 2 + 3
```

### Calculations

MoLang supports these arithmetic operators:

```text
+  addition
-  subtraction
*  multiplication
/  division
```

Parentheses can be used to group calculations:

```molang
dikhao (10 + 5) * 2
```

## How It Works

MoLang is implemented in three small stages:

1. **Lexer**: converts source code into tokens such as keywords, identifiers,
	numbers, operators, and parentheses.
2. **Parser**: checks the token sequence and builds the structure of the
	program.
3. **Execution environment**: stores variables and evaluates statements.

The lexer currently defines token types for `yaha`, `dikhao`, identifiers,
numbers, assignment, arithmetic operators, parentheses, new lines, and end of
file.

## Run the Project

Make sure Go is installed, then run:

```bash
go run main.go
```

At the current stage, the program prints a welcome message while the lexer,
parser, and execution environment are being implemented.

## Project Goal

MoLang is an educational project for understanding how programming languages
work internally, from reading source text to executing instructions.