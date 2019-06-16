package codewriter

import (
	"os"
	"strconv"
	"../parser"
)

type CodeWriter struct {
	fp *os.File
	parsingFileName string
	labelFirstCount int
	labelSecondCount int
}

func (cw *CodeWriter) Init(fp *os.File) {
	cw.fp = fp
	lines := getInitAssembly()
	cw.fp.WriteString(lines)
	cw.labelFirstCount = 0
	cw.labelSecondCount = 1
}

func (cw *CodeWriter) SetFileName(parsingFileName string) {
	cw.parsingFileName = parsingFileName
}

func (cw *CodeWriter) WriteArithmetic(command string) {
	var stackOperationAssemblyTable = map[string]string{
		"add": "M=D+M",
		"sub": "M=M-D",
		"neg": "M=-M",
		"eq": "D;JEQ",
		"gt": "D;JGT",
		"lt": "D;JLT",
		"and": "M=D&M",
		"or": "M=D|M",
		"not": "M=!M",
	}
	var lines string
	switch {
	case command == "add" || command == "sub" || command == "and" || command == "or":
		lines = getPopBinaryAssembly()
		lines += (stackOperationAssemblyTable[command] + "\n")
		lines += getAdvanceStackPointerAssembly()
	case command == "neg" || command == "not":
		lines = getPopUnaryAssembly()
		lines += (stackOperationAssemblyTable[command] + "\n")
		lines += getAdvanceStackPointerAssembly()
	case command == "eq" || command == "gt" || command == "lt":
		lines = getPopBinaryAssembly()
		lines += `D=M-D
		`
		lines += ("@LABEL" + strconv.Itoa(cw.labelFirstCount) + "\n")
		lines += (stackOperationAssemblyTable[command] + "\n")
		lines += `@SP
		A=M
		M=0
		`
		lines += ("@LABEL" + strconv.Itoa(cw.labelSecondCount) + "\n")
		lines += "0;JMP\n"
		lines += ("(LABEL" + strconv.Itoa(cw.labelFirstCount) + ")\n")
		lines += `@SP
		A=M
		M=-1
		`
		lines += ("(LABEL" + strconv.Itoa(cw.labelSecondCount) + ")\n")
		lines += getAdvanceStackPointerAssembly()

		cw.labelFirstCount += 2
		cw.labelSecondCount += 2
	}
	cw.fp.WriteString(lines)
}

func (cw *CodeWriter) WritePushPop(command int, segment string, index int) {
	var lines string
	if command == parser.C_PUSH {
		switch {
		case segment == "constant":
			lines = ("@" + strconv.Itoa(index))
			lines += `
			D=A
			@SP
			A=M
			M=D
			`
			lines += getAdvanceStackPointerAssembly()
		}
	}
	cw.fp.WriteString(lines)
}

func getInitAssembly() string {
	return `@256
	        D=A
	        @SP
	        M=D
	`
}

func getAdvanceStackPointerAssembly() string {
	return `@SP
	          M=M+1
	`
}

func getPopBinaryAssembly() string {
	return `@SP
			  AM=M-1
			  D=M
			  @SP
			  AM=M-1
	`
}

func getPopUnaryAssembly() string {
	return `@SP
	AM=M-1
	`
}
