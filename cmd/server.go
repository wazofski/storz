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
	mhr := store.New(sch, react.MetaHHandlerFactory(mem))
	rct := store.New(sch, react.ReactFactory(mhr))
	ssr := store.New(sch, react.StatusStripperFactory(rct))

	srv := rest.Server(sch, ssr)
	srv.Listen(8000)
}
