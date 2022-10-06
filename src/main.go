package main

import (
	"os"
)

func main() {
	if len(os.Args) != 3 {
		panic("Invalid argument")
	}
	bytes := OpenFile(os.Args[1])
	instructions := Parse(bytes)
	WriteFile(os.Args[2], ToMachineCode(instructions))
}
