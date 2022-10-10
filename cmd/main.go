package main

import (
	"fmt"
	"log"
	"os"

	"github.com/wazofski/store/mgen"
)

func main() {
	argLength := len(os.Args[1:])
	if argLength < 2 {
		log.Fatalln("missing arguments: go run builder.go model destination")
	}

	fmt.Println(mgen.Generate(os.Args[1], os.Args[2]))
}
