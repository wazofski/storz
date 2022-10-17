echo ">>>>>>>>>>>>> STARTING COMMON SUITE"

go clean -cache

go test -ginkgo.v -args store=0
go test -ginkgo.v -args store=1
go test -ginkgo.v -args store=3

trap "kill 0" EXIT

cd ..
go run cmd/server.go &
cd commontests

sleep 1

go test -ginkgo.v -args store=2

