echo "########### CLEARING GO CACHE ###############"
go clean -cache

echo "########### STARTING TEST SUITE #############"
ginkgo -r -focus "storz"
ginkgo -r -focus "mgen"

ginkgo -r -focus "react"
ginkgo -r -focus "client"

cd test
./run_tests.sh
cd ..
