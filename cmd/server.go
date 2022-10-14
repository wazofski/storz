package main

import (
	"github.com/wazofski/store"
	"github.com/wazofski/store/generated"
	"github.com/wazofski/store/memory"
	"github.com/wazofski/store/rest"
)

func main() {
	sch := generated.Schema()
	mem := store.New(sch, memory.Factory())

	srv := rest.Server(sch, mem)
	srv.Listen(8000)
}
