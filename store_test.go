package store_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/store/mgen"
)

var _ = Describe("store", func() {
	It("mgen can generate", func() {
		err := mgen.Generate("mgen/testmodel")
		Expect(err).To(BeNil())
	})
})
