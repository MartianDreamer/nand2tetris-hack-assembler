package main

import "os"

func OpenFile(path string) []byte  {
	file, err := os.ReadFile(path)
	if (err != nil) {
		panic("Failed to open file")
	}
	return file
} 