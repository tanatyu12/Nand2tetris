package main

import (
	"./parser"
	"./code"
	"./symboltable"
	"fmt"
	"strconv"
	"flag"
	"strings"
	"os"
)

func main() {
	// get assembly file
	flag.Parse()
	asmFileName := flag.Arg(0)

	// open assembly file (for read)
	asmFp, err := os.Open(asmFileName)
	if err != nil {
		return
	}
	defer asmFp.Close()

	// open binary file (for write)
	hackFileName := strings.Replace(asmFileName, ".asm", ".hack", 1)
	hackFp, err := os.Create(hackFileName)
	if err != nil {
		return
	}
	defer hackFp.Close()

	st := new(symboltable.SymbolTable)
	st.Init()

	p1 := new(parser.Parser)
	p1.Init(asmFp)
	programCounter := -1
	for p1.HasMoreCommand() {
		p1.Advance()

		symbol := p1.Symbol()
		if symbol == "" {
			break
		}
		commandType := p1.CommandType()

		if commandType == parser.L_COMMAND {
			symbol = strings.TrimLeft(symbol, "(")
			symbol = strings.TrimRight(symbol, ")")
			st.AddEntry(symbol, programCounter + 1)
		} else {
			programCounter++
		}
	}

	p2 := new(parser.Parser)
	asmFp.Seek(0, 0)
	p2.Init(asmFp)
	variableAddress := 16
	for p2.HasMoreCommand() {
		p2.Advance()

		if p2.Symbol() == "" {
			break
		}
		commandType := p2.CommandType()

		if commandType == parser.A_COMMAND {
			symbol := p2.Symbol()
			if i, err := strconv.Atoi(symbol); err == nil {
				symbolBinary := fmt.Sprintf("%016b", i)
				hackFp.WriteString(symbolBinary + "\n")
			} else {
				if !st.Contains(symbol) {
					st.AddEntry(symbol, variableAddress)
					variableAddress++
				}
				symbolBinary := fmt.Sprintf("%016b", st.GetAddress(symbol))
				hackFp.WriteString(symbolBinary + "\n")
			}
		}

		if commandType == parser.C_COMMAND {
			compMnemonic := p2.Comp()
			destMnemonic := p2.Dest()
			jumpMnemonic := p2.Jump()

			compBinary := code.Comp(compMnemonic)
			destBinary := code.Dest(destMnemonic)
			jumpBinary := code.Jump(jumpMnemonic)
			hackFp.WriteString("111" + compBinary + destBinary + jumpBinary + "\n")
		}
	}
}