package store_test

import (
	"log"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/store"
	"github.com/wazofski/store/generated"
)

var _ = Describe("mgen", func() {
	Describe("class", func() {
		It("factory", func() {
			world := generated.WorldFactory()
			Expect(world).ToNot(BeNil())
			Expect(world.Spec()).ToNot(BeNil())
			Expect(world.Status()).ToNot(BeNil())
		})

		It("setters and getters ", func() {
			world := generated.WorldFactory()
			world.Spec().SetName("abc")
			Expect(world.Spec().Name()).To(Equal("abc"))

			world.Spec().Nested().SetAlive(true)
			Expect(world.Spec().Nested().Alive()).To(BeTrue())

			world.Spec().Nested().SetCounter(10)
			Expect(world.Spec().Nested().Counter()).To(Equal(10))

			world.Status().SetDescription("qwe")
			Expect(world.Status().Description()).To(Equal("qwe"))
		})

		It("metadata", func() {
			world := generated.WorldFactory()
			Expect(world.Metadata().Kind()).To(Equal(store.ObjectKind("World")))
		})

		It("serialization", func() {
			world := generated.WorldFactory()
			world.Spec().Nested().SetCounter(10)
			world.Spec().Nested().SetAlive(true)
			world.Spec().SetName("abc")
			world.Status().SetDescription("qwe")

			log.Println(string(world.Serialize()))
		})

	})

})
