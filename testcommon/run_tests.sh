echo ">>>>>>>>>>>>> STARTING COMMON SUITE"

go test -ginkgo.v -args store=0
go test -ginkgo.v -args store=1

trap "kill 0" EXIT

go clean -cache

cd ..
go run cmd/server.go &
cd testcommon

sleep 1

go test -ginkgo.v -args store=2

