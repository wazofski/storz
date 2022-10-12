//go:build ignore

package main

import (
	"log"

	"github.com/wazofski/store/mgen"
)

func main() {
	err := mgen.Generate()
	if err != nil {
		log.Fatalln(err)
	}
}
