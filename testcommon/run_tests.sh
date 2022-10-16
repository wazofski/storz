echo ">>>>>>>>>>>>> STARTING COMMON SUITE"

go clean -cache

go test -ginkgo.v -args store=0
go test -ginkgo.v -args store=1

trap "kill 0" EXIT

cd ..
go run cmd/server.go &
cd testcommon

sleep 1

go test -ginkgo.v -args store=2

