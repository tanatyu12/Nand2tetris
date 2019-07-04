package parser

import (
	"bufio"
	"os"
	"strings"
	// "fmt"

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
		} else {
			p.currentCommand = " "
		}
	} else {
		p.currentCommand = line
	}
}

const (
	A_COMMAND = iota
	C_COMMAND
	L_COMMAND
)

func (p *Parser) CommandType() int {
	isTypeA := strings.HasPrefix(p.currentCommand, "@")
	isTypeC := strings.Contains(p.currentCommand, "=") || strings.Contains(p.currentCommand, ";")
	if isTypeA {
		return A_COMMAND
	} else if isTypeC {
		return C_COMMAND
	} else {
		return L_COMMAND
	}
}

func (p *Parser) Symbol() string {
	currentCommand := p.currentCommand
	if p.CommandType() == A_COMMAND {
		currentCommand = strings.TrimLeft(currentCommand, "@")
	}
	return currentCommand
}

func (p *Parser) Dest() string {
	if strings.Contains(p.currentCommand, "=") {
		slice := strings.Split(p.currentCommand, "=")
		return slice[0]
	}
	return ""
}

func (p *Parser) Comp() string {
	if strings.Contains(p.currentCommand, "=") {
		slice := strings.Split(p.currentCommand, "=")
		return slice[1]
	} else {
		slice := strings.Split(p.currentCommand, ";")
		return slice[0]
	}
}

func (p *Parser) Jump() string {
	if strings.Contains(p.currentCommand, ";") {
		slice := strings.Split(p.currentCommand, ";")
		return slice[1]
	}
	return ""	
}