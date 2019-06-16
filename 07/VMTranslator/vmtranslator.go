package main

import (
	"flag"
	"strings"
	"io/ioutil"
	"path/filepath"
	"./parser"
	"./codewriter"
	"os"
)

func main() {
	flag.Parse()
	arg := flag.Arg(0)

	var asmFileName string
	var vmFilePaths []string
	if strings.HasSuffix(arg, ".vm") {
		vmFilePaths = append(vmFilePaths, arg)
		asmFileName = strings.Replace(arg, ".vm", ".asm", 1)
	} else {
		asmFileName = arg + ".asm"
		vmFilePaths = getVMFilePaths(arg)
	}

	asmFp, err := os.Create(asmFileName)
	if err != nil {
		panic(err)
	}
	defer asmFp.Close()

	c := new(codewriter.CodeWriter)
	c.Init(asmFp)

	for _, vmFilePath := range vmFilePaths {
		c.SetFileName(vmFilePath)
		vmFp, err := os.Open(vmFilePath)
		if err != nil {
			panic(err)
		}

		p := new(parser.Parser)
		p.Init(vmFp)
		for p.HasMoreCommand() {
			p.Advance()
			commandType := p.CommandType()
			arg1 := p.Arg1()
			arg2 := 0
			if commandType == parser.C_PUSH || commandType == parser.C_POP || commandType == parser.C_FUNCTION || commandType == parser.C_CALL {
				arg2 = p.Arg2()
			}
			if commandType == parser.C_PUSH || commandType == parser.C_POP {
				c.WritePushPop(commandType, arg1, arg2)
			}
			if commandType == parser.C_ARITHMETIC {
				c.WriteArithmetic(arg1)
			}
		}
		vmFp.Close()
	}
}

func getVMFilePaths(dirPath string) []string {
	fileInfoList, err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	vmFilePaths := []string{}
	for _, fileInfo := range fileInfoList {
		currentPath := filepath.Join(dirPath, fileInfo.Name())
		if fileInfo.IsDir() {
			vmFilePaths = append(vmFilePaths, getVMFilePaths(currentPath)...)
			continue
		} else {
			if !strings.HasSuffix(currentPath, ".vm") {
				continue
			}
			vmFilePaths = append(vmFilePaths, currentPath)
		}
	}

	return vmFilePaths
}