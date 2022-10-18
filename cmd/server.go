package main

import (
	"github.com/wazofski/storz/generated"
	"github.com/wazofski/storz/memory"
	"github.com/wazofski/storz/react"
	"github.com/wazofski/storz/rest"
	"github.com/wazofski/storz/store"
)

func main() {
	sch := generated.Schema()
	mem := store.New(sch, memory.Factory())
	mhr := store.New(sch, react.MetaHHandlerFactory(mem))
	rct := store.New(sch, react.ReactFactory(mhr))
	ssr := store.New(sch, react.StatusStripperFactory(rct))

	srv := rest.Server(sch, ssr)
	srv.Listen(8000)
}
