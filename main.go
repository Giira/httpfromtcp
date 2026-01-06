package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	for line := range getLinesChannel(f) {
		fmt.Printf("read: %v\n", line)
	}
}
