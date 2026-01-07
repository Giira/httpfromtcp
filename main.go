package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	f, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("error listening on port 42069: %v", err)
	}
	defer f.Close()

	for {
		con, err := f.Accept()
		if err != nil {
			log.Fatalf("error accepting connection: %v", err)
		}
		fmt.Println("Connection accepted")
		for line := range getLinesChannel(con) {
			fmt.Println(line)
		}
	}
}
