package main

import (
	"fmt"
	"io"
	"log"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	chanOut := make(chan string)
	var currentString string

	for {
		b := make([]byte, 8)
		_, err := f.Read(b)
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
			chanOut <- currentString
			currentString = sSlice[1]
		}
	}
	return chanOut
}
