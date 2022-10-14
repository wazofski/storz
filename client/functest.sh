trap "kill 0" EXIT

go clean -cache

cd ..
go run cmd/server.go > sample.out &
sleep 1

ginkgo -r -focus "client" -v
# go test -ginkgo.v -ginkgo.focus "client"
