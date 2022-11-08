package main

// func main() {
// 	sch := generated.Schema()

// 	// mem := store.New(sch, mongo.Factory("mongodb://localhost:27017/", "storz"))
// 	mem := store.New(sch, memory.Factory())
// 	mhr := store.New(sch, react.MetaHHandlerFactory(mem))
// 	rct := store.New(sch, react.ReactFactory(mhr))
// 	ssr := store.New(sch, react.StatusStripperFactory(rct))

// 	// srv := browser.Server(sch, ssr)
// 	srv := rest.Server(sch, ssr)
// 	srv.Listen(8000)
// }
