package codewriter

import (
	"os"
	"strconv"
	"../parser"
)

type CodeWriter struct {
	fp *os.File
	fileName string
	labelFirstCount int
	labelSecondCount int
}

func (cw *CodeWriter) Init(fp *os.File) {
	cw.fp = fp
	// lines := getInitAsm()
	// cw.fp.WriteString(lines)
	cw.labelFirstCount = 0
	cw.labelSecondCount = 1
}

func (cw *CodeWriter) SetFileName(fileName string) {
	cw.fileName = fileName
}

func (cw *CodeWriter) WriteArithmetic(command string) {
	var stackOperationAsmTable = map[string]string{
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
			lines = getBinaryBackStackPointerAsm()
			lines += (stackOperationAsmTable[command] + "\n")
			lines += getAdvanceStackPointerAsm()
	    case command == "neg" || command == "not":
			lines = getUnaryBackStackPointerAsm()
			lines += (stackOperationAsmTable[command] + "\n")
			lines += getAdvanceStackPointerAsm()
	    case command == "eq" || command == "gt" || command == "lt":
			lines = getBinaryBackStackPointerAsm()
			lines += "D=M-D\n"
			lines += ("@LABEL" + strconv.Itoa(cw.labelFirstCount) + "\n")
			lines += (stackOperationAsmTable[command] + "\n")
			lines += "@SP\nA=M\nM=0\n"
			lines += ("@LABEL" + strconv.Itoa(cw.labelSecondCount) + "\n")
			lines += "0;JMP\n"
			lines += ("(LABEL" + strconv.Itoa(cw.labelFirstCount) + ")\n")
			lines += "@SP\nA=M\nM=-1\n"
			lines += ("(LABEL" + strconv.Itoa(cw.labelSecondCount) + ")\n")
			lines += getAdvanceStackPointerAsm()

			cw.labelFirstCount += 2
			cw.labelSecondCount += 2
	}
	cw.fp.WriteString(lines)
}

var registerNameMap = map[string]string{
	"local": "LCL",
	"argument": "ARG",
	"this": "THIS",
	"that": "THAT",
}

func (cw *CodeWriter) WritePushPop(command int, segment string, index int) {
	var lines = ""
	switch command {
		case parser.C_PUSH:
			lines += getPushAsm(segment, index, cw.fileName)
		case parser.C_POP:
			lines += getPopAsm(segment, index, cw.fileName)
	}
	cw.fp.WriteString(lines)
}


func getInitAsm() string {
	// set base addresses into segment register
	baseAddresses := map[string]int{
		"SP": 256,
		"LCL": 24576,
		"ARG": 28160,
	}
	lines := ""
	for segment, baseAddress := range baseAddresses {
		lines += ("@" + strconv.Itoa(baseAddress) + "\n")
		lines += "D=A\n"
		lines += ("@" + segment + "\n")
		lines += "M=D\n"
	}
	return lines
}

func getPushAsm(segment string, index int, fileName string) string {
	var lines = ""
	switch {
		case segment == "constant":
			lines += ("@" + strconv.Itoa(index) + "\n")
			lines += "D=A\n@SP\nA=M\nM=D\n"
			lines += getAdvanceStackPointerAsm()
		case segment == "local" || segment == "argument" || segment == "this" || segment == "that":
			// read source value
			lines += ("@" + strconv.Itoa(index) + "\n")
			lines += "D=A\n@" + registerNameMap[segment] + "\nA=D+M\nD=M\n"
			// push
			lines += "@SP\nA=M\nM=D\n"
			lines += getAdvanceStackPointerAsm()
		case segment == "pointer" || segment == "temp":
			var baseAddress int
			if segment == "pointer" {
				baseAddress = 3
			} else {
				baseAddress = 5
			}
			lines += ("@" + strconv.Itoa(baseAddress + index) + "\n")
			lines += "D=M\n@SP\nA=M\nM=D\n"
			lines += getAdvanceStackPointerAsm()
		case segment == "static":
			lines += ("@" + fileName + "." + strconv.Itoa(index) + "\n")
			lines += "D=M\n@SP\nA=M\nM=D\n"
			lines += getAdvanceStackPointerAsm()
	}
	return lines
}

func getPopAsm(segment string, index int, fileName string) string {
	var lines = ""
	switch {
		case segment == "local" || segment == "argument" || segment == "this" || segment == "that":
			// move dest address to RAM[13]
			lines += ("@" + strconv.Itoa(index) + "\n")
			lines += "D=A\n@" + registerNameMap[segment] + "\nD=D+M\n@R13\nM=D\n"
			// pop
			lines += getUnaryBackStackPointerAsm()
			lines += "D=M\n@R13\nA=M\nM=D\n"
		case segment == "pointer" || segment == "temp":
			lines += getUnaryBackStackPointerAsm()
			lines += "D=M\n"
			var baseAddress int
			if segment == "pointer" {
				baseAddress = 3
			} else {
				baseAddress = 5
			}
			lines += ("@" + strconv.Itoa(baseAddress + index) + "\n")
			lines += "M=D\n"
		case segment == "static":
			lines += getUnaryBackStackPointerAsm()
			lines += "D=M\n"
			lines += ("@" + fileName + "." + strconv.Itoa(index) + "\n")
			lines += "M=D\n"
	}	
	return lines
}

func getAdvanceStackPointerAsm() string {
	return "@SP\nM=M+1\n"
}

func getBinaryBackStackPointerAsm() string {
	return "@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\n"
}

func getUnaryBackStackPointerAsm() string {
	return "@SP\nAM=M-1\n"
}
