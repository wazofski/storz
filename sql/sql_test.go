package sql_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/store/generated"
)

var _ = Describe("sql", func() {

	It("can initialize a db", func() {
		Expect(stc).ToNot(BeNil())

		ret, err := stc.List(ctx, generated.SecondWorldIdentity(""))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
	})

})
