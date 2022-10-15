echo ">>>>>>>>> CLIENT TESTS"

trap "kill 0" EXIT

go clean -cache

cd ..
go run cmd/server.go &
sleep 1

ginkgo -r -focus "client" -v
# go test -ginkgo.v -ginkgo.focus "client"
