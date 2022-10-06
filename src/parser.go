package main

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func exploreSymbol(fileContent string) []string {
	nixSep := "\n"
	windowsSep := "\r\n"
	var lines []string
	if runtime.GOOS != "windows" {
		lines = strings.Split(fileContent, nixSep)
	} else {
		lines = strings.Split(fileContent, windowsSep)
	}
	for i := 0; i < len(lines); i++ {
		lines[i] = strings.ReplaceAll(lines[i], " ", "")
		if lines[i] == "" || strings.HasPrefix(lines[i], "//") {
			lines = append(lines[:i], lines[i+1:]...)
			i -= 1
		} else if strings.Contains(lines[i], "//") {
			index := strings.Index(lines[i], "//")
			lines[i] = lines[i][:index]
		}
	}
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "(") && strings.HasSuffix(lines[i], ")") {
			symbol := lines[i][1 : len(lines[i])-1]
			if _, ok := SymbolMap[symbol]; ok {
				panic("re-declaring is not allowed - line " + strconv.Itoa(i))
			}
			SymbolMap[symbol] = i
			lines = append(lines[:i], lines[i+1:]...)
			i -= 1
		}
	}
	return lines
}

func Parse(bytes []byte) []Instruction {
	byteString := string(bytes)
	trimmedSpaceString := exploreSymbol(byteString)
	instructions := make([]Instruction, len(trimmedSpaceString))
	for i, line := range trimmedSpaceString {
		ins, err := parseLine(line)
		if err != nil {
			panic("Invalid instruction " + strconv.Itoa(i))
		}
		instructions[i] = ins
	}
	return instructions
}

func parseLine(assemblyInstruction string) (Instruction, error) {
	if assemblyInstruction[0] == '@' {
		return parseAInstruction(assemblyInstruction), nil
	}
	return parseCInstruction(assemblyInstruction)

}

func parseAInstruction(assemblyInstruction string) Instruction {
	if address, err := strconv.Atoi(assemblyInstruction[1:]); err == nil {
		return Instruction(address)
	}
	if address, ok := SymbolMap[assemblyInstruction[1:]]; ok {
		return Instruction(address)
	}
	address := AvailableRamPos
	SymbolMap[assemblyInstruction[1:]] = address
	AvailableRamPos++
	return Instruction(address)
}

func parseCInstruction(assemblyInstruction string) (Instruction, error) {
	compAndJump := strings.Split(assemblyInstruction, ";")
	err := errors.New("invalid instruction")
	var (
		compPhrase  string
		jumpPhrase  string
		computation string
		destination string
	)
	if len(compAndJump) == 1 {
		compPhrase = compAndJump[0]
	} else if len(compAndJump) == 2 {
		compPhrase = compAndJump[0]
		jumpPhrase = compAndJump[1]
	} else {
		return Instruction(-1), err
	}
	destAndComp := strings.Split(compPhrase, "=")
	if len(destAndComp) == 1 {
		computation = destAndComp[0]
	} else if len(destAndComp) == 2 {
		destination = destAndComp[0]
		computation = destAndComp[1]
	} else {
		return Instruction(-1), err
	}
	jumpInt := JumpMap[jumpPhrase]
	compInt := ComputationMap[computation]
	destInt := DestinationMap[destination]
	instruction := jumpInt + destInt<<3 + compInt<<6 + 0b111<<13
	return Instruction(instruction), nil
}

func ToMachineCode(instructions []Instruction) []string {
	lines := make([]string, len(instructions))
	for i, ins := range instructions {
		lines[i] = fmt.Sprintf("%016b\n", int(ins))
	}
	return lines
}
