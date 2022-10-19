package main

import (
	"github.com/wazofski/storz/browser"
	"github.com/wazofski/storz/generated"
	"github.com/wazofski/storz/mongo"
	"github.com/wazofski/storz/react"
	"github.com/wazofski/storz/store"
)

func main() {
	sch := generated.Schema()

	// st := store.New(
	// 	sch, mongo.Factory("mongodb://localhost:27017/", "storz"))

	// obj := generated.WorldFactory()
	// obj.Spec().SetName("ajkshsjkh")

	// ctx := context.Background()

	// ret, err := st.Create(ctx, obj)
	// if err != nil {
	// 	log.Println(err)
	// }

	// if ret != nil {
	// 	_, err = st.Get(ctx, ret.Metadata().Identity())
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }

	// _, err = st.Get(ctx, generated.WorldIdentity("ajkshsjkh"))

	// if err != nil {
	// 	log.Println(err)
	// }

	mem := store.New(sch, mongo.Factory("mongodb://localhost:27017/", "storz"))
	// mem := store.New(sch, memory.Factory())
	mhr := store.New(sch, react.MetaHHandlerFactory(mem))
	rct := store.New(sch, react.ReactFactory(mhr))
	ssr := store.New(sch, react.StatusStripperFactory(rct))

	srv := browser.Server(sch, ssr)
	// srv := rest.Server(sch, ssr)
	srv.Listen(8000)
}
