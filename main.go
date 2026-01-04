package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	var addr int64
	var currentString string

	for {
		b := make([]byte, 8)
		_, err := f.Seek(addr, io.SeekStart)
		if err != nil {
			log.Fatalf("error finding next set of bytes: %v", err)
		}
		_, err1 := f.Read(b)
		if err1 == io.EOF {
			fmt.Printf("read: %v\n", currentString)
			break
		}

		s := string(b)
		sSlice := strings.Split(s, "\n")
		currentString += sSlice[0]
		if len(sSlice) != 1 {
			fmt.Printf("read: %v\n", currentString)
			currentString = sSlice[1]
		}
		addr += 8
	}
}
