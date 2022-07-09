package main

import (
	"fmt"
	"log"
	"os"

	proto_reader "github.com/Not-cottage-cheese-but-cottage-cheese/final-go/proto"
)

func main() {
	if len(os.Args) != 3 {
		log.Panicln("invalid arguments count")
	}

	protoDirPath, binaryDirPath := os.Args[1], os.Args[2]
	fmt.Println(protoDirPath, binaryDirPath)

	proto_reader.GetDescriptors(protoDirPath)
}
