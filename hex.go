package main

import (
	"fmt"
	"os"
	"strings"

	mmap "github.com/edsrzf/mmap-go"
)

// dumpFile memory maps the file given as input and processes it extracting
// the address, hex and ascii string representations of the file
func dumpFile(fileName string) (addresses, hex, text []string) {
	// Open the file
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	// Memory map the file
	mappedFile, err := mmap.Map(f, mmap.RDONLY, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer mappedFile.Unmap()

	addresses = make([]string, 0, len(mappedFile)/16)
	hex = make([]string, 0, len(mappedFile)/16)
	text = make([]string, 0, len(mappedFile)/16)

	// Process the memory mapped file 16 bytes at a time
	for i := 0; i < len(mappedFile); i += 16 {
		offset := i + 16
		if offset > len(mappedFile) {
			offset = len(mappedFile)
		}
		addresses = append(addresses, fmt.Sprintf("%08x", i))
		hex = append(hex, fmt.Sprintf("%x", mappedFile[i:offset]))
		str := fmt.Sprintf("%s", mappedFile[i:offset])
		str = strings.ReplaceAll(str, "\n", "\\n")
		text = append(text, str)
	}

	return
}
