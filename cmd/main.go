package main

import (
	"log"
	"os"

	"github.com/wazofski/store/mgen"
)

func main() {
	argLength := len(os.Args[1:])
	if argLength < 2 {
		log.Fatalln("missing arguments: go run builder.go model destination")
	}

	mgen.Generate(os.Args[1], os.Args[2])
}
