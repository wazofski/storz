ginkgo -r -focus "store"
ginkgo -r -focus "mgen"
ginkgo -r -focus "memory"
ginkgo -r -focus "react"

cd client
./run_tests.sh
cd ..

# cd testcommon
# ./run_tests.sh
# cd ..

