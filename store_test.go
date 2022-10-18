package store_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/storz/mgen"
)

var _ = Describe("store", func() {
	It("mgen can generate", func() {
		err := mgen.Generate("testmodel")
		Expect(err).To(BeNil())
	})
})
