package store_test

import (
	"encoding/json"
	"log"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/store"
	"github.com/wazofski/store/generated"
)

var _ = Describe("mgen", func() {
	Describe("class", func() {
		It("can factory", func() {
			world := generated.WorldFactory()
			Expect(world).ToNot(BeNil())
			Expect(world.Spec()).ToNot(BeNil())
			Expect(world.Status()).ToNot(BeNil())
		})

		It("can call setters and getters ", func() {
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

		It("has metadata", func() {
			world := generated.WorldFactory()
			Expect(world.Metadata().Kind()).To(Equal(store.ObjectKind("World")))
		})

		It("can deserialize", func() {
			world := generated.WorldFactory()
			world.Spec().Nested().SetCounter(10)
			world.Spec().Nested().SetAlive(true)
			world.Spec().SetName("abc")
			world.Status().SetDescription("qwe")
			world.Status().SetList([]generated.NestedWorld{
				generated.NestedWorldFactory(),
				generated.NestedWorldFactory(),
			})

			world.Status().SetMap(map[string]generated.NestedWorld{
				"a": generated.NestedWorldFactory(),
				"b": generated.NestedWorldFactory(),
			})

			world.Status().Map()["a"].SetL1([]bool{false, false, true})

			data, err := json.MarshalIndent(world, "", "  ")
			Expect(err).To(BeNil())

			log.Println(string(data))

			newWorld := generated.WorldFactory()
			err = json.Unmarshal(data, &newWorld)

			Expect(err).To(BeNil())
			Expect(newWorld.Spec().Nested().Alive()).To(BeTrue())
			Expect(newWorld.Spec().Nested().Counter()).To(Equal(10))
			Expect(newWorld.Spec().Name()).To(Equal("abc"))
			Expect(newWorld.Status().Description()).To(Equal("qwe"))
			Expect(len(newWorld.Status().List())).To(Equal(2))

			data2, err := json.MarshalIndent(newWorld, "", "  ")
			Expect(err).To(BeNil())
			Expect(data).To(Equal(data2))
		})

		It("has working schema", func() {
			world := generated.WorldFactory()
			obj := generated.ObjectForKind(string(world.Metadata().Kind()))
			Expect(obj).ToNot(BeNil())
			anotherWorld := obj.(generated.World)
			Expect(anotherWorld).ToNot(BeNil())
		})

		It("can clone objects", func() {
			world := generated.WorldFactory()
			world.Spec().Nested().SetCounter(10)
			world.Spec().Nested().SetAlive(true)
			world.Spec().SetName("abc")
			world.Status().SetDescription("qwe")
			world.Status().SetList([]generated.NestedWorld{
				generated.NestedWorldFactory(),
				generated.NestedWorldFactory(),
			})

			world.Status().SetMap(map[string]generated.NestedWorld{
				"a": generated.NestedWorldFactory(),
				"b": generated.NestedWorldFactory(),
			})

			world.Status().Map()["a"].SetL1([]bool{false, false, true})

			newWorld := world.Clone().(generated.World)
			Expect(newWorld.Spec().Nested().Alive()).To(BeTrue())
			Expect(newWorld.Spec().Nested().Counter()).To(Equal(10))
			Expect(newWorld.Spec().Name()).To(Equal("abc"))
			Expect(newWorld.Status().Description()).To(Equal("qwe"))
			Expect(len(newWorld.Status().List())).To(Equal(2))
		})

	})

})
