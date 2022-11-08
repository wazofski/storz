echo ">>>>>>>>>>>>> STARTING COMMON SUITE"

go clean -cache

go test -ginkgo.v -args store=memory
go test -ginkgo.v -args store=client
go test -ginkgo.v -args store=react
go test -ginkgo.v -args store=sqlite
go test -ginkgo.v -args store=mysql
go test -ginkgo.v -args store=mongo
