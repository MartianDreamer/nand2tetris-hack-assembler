package main

import "os"

func OpenFile(path string) []byte {
	file, err := os.ReadFile(path)
	if err != nil {
		panic("Failed to open file")
	}
	return file
}

func WriteFile(path string, content []string) {
	file, err := os.Create(path)
	if err != nil {
		panic("failed to write file")
	}
	for _, line := range content {
		file.WriteString(line)
	}
}
