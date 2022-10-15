ginkgo -r -focus "store"
ginkgo -r -focus "mgen"
ginkgo -r -focus "memory"
ginkgo -r -focus "react"

ginkgo -r -focus "negative"

cd client
./run_tests.sh
cd ..

cd testcommon
./run_tests.sh
cd ..

