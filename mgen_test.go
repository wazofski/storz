package store_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/store/generated"
)

var _ = Describe("mgen", func() {
	It("class factory", func() {
		world := generated.WorldFactory()
		Expect(world).ToNot(BeNil())
		Expect(world.Spec()).ToNot(BeNil())
		Expect(world.Status()).ToNot(BeNil())
	})

	It("class setters and getters ", func() {
		world := generated.WorldFactory()
		world.Spec().SetName("abc")
		Expect(world.Spec().Name()).To(Equal("abc"))

		world.Spec().Nested().SetAlive(true)
		Expect(world.Spec().Nested().Alive()).To(BeTrue())

		world.Spec().Nested().SetCounter(10)
		Expect(world.Spec().Nested().Counter()).To(Equal(10))
	})

})
