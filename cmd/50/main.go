package main

import (
	"log"
	"os"

	proto_reader "github.com/Not-cottage-cheese-but-cottage-cheese/final-go/proto"
)

func main() {
	if len(os.Args) != 3 {
		log.Panicln("invalid arguments count")
	}

	protoDirPath, _ := os.Args[1], os.Args[2]

	proto_reader.GetDescriptors(protoDirPath)
}
