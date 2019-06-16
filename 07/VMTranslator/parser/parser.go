package parser

import (
	"bufio"
	"os"
	"strings"
	"strconv"
)

type Parser struct {
	currentCommand string
	sc *bufio.Scanner
}

func (p *Parser) Init(fp *os.File) {
	p.sc = bufio.NewScanner(fp)
}

func (p *Parser) HasMoreCommand() bool {
	return p.sc.Scan()
}

func (p *Parser) Advance() {
	line := p.sc.Text()
	if commentIdx := strings.Index(line, "//"); commentIdx > -1 {
		line = line[:commentIdx]
	}
	line = strings.TrimSpace(line)
	if line == "" {
		if p.HasMoreCommand() {
			p.Advance()
		}
	} else {
		p.currentCommand = line
	}
}

const (
	C_ARITHMETIC = iota
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)

func (p *Parser) CommandType() int {
	commandSplit := strings.Split(p.currentCommand, " ")
	opecode := commandSplit[0]

	arithmeticCommands := []string{
		"add",
		"sub",
		"neg",
		"eq",
		"gt",
		"lt",
		"and",
		"or",
		"not",
	}
	for _, arithmeticCommand := range arithmeticCommands {
		if opecode == arithmeticCommand {
			return C_ARITHMETIC
		}
	}

	if opecode == "push" {
		return C_PUSH
	}
	return C_ARITHMETIC
}

func (p *Parser) Arg1() string {
	if p.CommandType() == C_ARITHMETIC {
		return p.currentCommand
	}
	arg1 := strings.Split(p.currentCommand, " ")[1]
	return arg1
}

func (p *Parser) Arg2() int {
	arg2 := strings.Split(p.currentCommand, " ")[2]
	arg2Int, _ := strconv.Atoi(arg2)
	return arg2Int
}