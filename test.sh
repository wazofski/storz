ginkgo -r -focus "store"
ginkgo -r -focus "mgen"
ginkgo -r -focus "memory"
ginkgo -r -focus "react"

cd client
./functest.sh
cd ..

