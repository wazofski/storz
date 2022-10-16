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
	mh := store.New(sch, react.MetaHHandlerFactory(mem))
	r := store.New(sch, react.Factory(mh))
	ss := store.New(sch, react.StatusStripperFactory(r))

	srv := rest.Server(sch, ss)
	srv.Listen(8000)
}
