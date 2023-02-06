package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	input := bufio.NewReader(os.Stdin)
	lineBytes, _, _ := input.ReadLine()
	line := string(lineBytes)
	lexer := NewLexer(line)
	parser := NewParser(lexer)
	program := parser.statement()
	fmt.Println(program)
}
