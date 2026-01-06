package main

import (
	"io"
	"log"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	chanOut := make(chan string)
	go func() {
		var currentString string
		defer close(chanOut)

		for {
			b := make([]byte, 8)
			_, err := f.Read(b)
			if err == io.EOF {
				if len(currentString) > 0 {
					chanOut <- currentString
					break
				}
			}
			if err != nil {
				log.Fatalf("error reading: %v\n", err)
			}
			s := string(b)
			sSlice := strings.Split(s, "\n")
			currentString += sSlice[0]
			if len(sSlice) != 1 {
				chanOut <- currentString
				currentString = sSlice[1]
			}
		}
	}()
	return chanOut
}
