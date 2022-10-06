package main

import (
	"strconv"
	"strings"
)

type InstructionType int

const (
	TypeA InstructionType = 0
	TypeC
)

type Instruction struct {
	Type     InstructionType
	Elements []string
}

func exploreSymbol(fileContent string) []string {
	nixSep := "\n"
	windowsSep := "\r\n"
	var lines []string
	if strings.Contains(fileContent, nixSep) {
		lines = strings.Split(fileContent, nixSep)
	} else {
		lines = strings.Split(fileContent, windowsSep)
	}
	for i := 0; i <= len(lines); i++ {
		lines[i] = strings.ReplaceAll(lines[i], " ", "")
		if lines[i] == "" {
			lines = append(lines[:i], lines[i+1:]...)
			i--
		}
	}
	for i := 0; i <= len(lines); i++ {
		if strings.HasPrefix(lines[i], "(") && strings.HasSuffix(lines[i], ")") {
			symbol := lines[1][1 : len(lines[1])-1]
			if _, ok := SymbolMap[symbol]; ok {
				panic("re-declaring is not allowed - line " + string(i))
			}
			SymbolMap[symbol] = int64(i)
			lines = append(lines[:i], lines[i+1:]...)
			i--
		}
	}
	return lines
}

func Parse(bytes []byte) (rs []Instruction) {
	return nil
}

func parseline(assemblyInstruction string) Instruction {
	if assemblyInstruction[0] == '@' {
		address, err := strconv.Atoi(assemblyInstruction[1:])
		if err != nil {
			panic("invalid instruction")
		}
		return Instruction{TypeA, []string{strconv.FormatInt(int64(address), 2)}}
	}
	return Instruction{}
}
