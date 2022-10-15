package common_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/store/generated"
)

var _ = Describe("common", func() {

	It("can POST objects", func() {
		w := generated.WorldFactory()

		w.Spec().SetName("abc")
		w.Spec().SetDescription("def")

		ret, err := clt.Create(ctx, w)

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret.Metadata().Identity())).ToNot(Equal(0))
	})

})
