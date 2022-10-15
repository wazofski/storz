package main

import (
	"github.com/wazofski/store"
	"github.com/wazofski/store/generated"
	"github.com/wazofski/store/memory"
	"github.com/wazofski/store/react"
	"github.com/wazofski/store/rest"
)

func main() {
	sch := generated.Schema()
	mem := store.New(sch, memory.Factory())
	reactor := store.New(sch, react.Factory(mem))

	srv := rest.Server(sch, reactor)
	srv.Listen(8000)
}
