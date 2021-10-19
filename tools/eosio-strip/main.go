package main

import "os"

func main() {
	inFile := os.Args[1]
	outFile := os.Args[2]
	wasmCheckSection(inFile, outFile)
}
