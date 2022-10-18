ginkgo -r -focus "storz"
ginkgo -r -focus "mgen"
ginkgo -r -focus "memory"
ginkgo -r -focus "react"

cd client
./run_tests.sh
cd ..

# cd common
# ./run_tests.sh
# cd ..

