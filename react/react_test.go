package react_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/store/generated"
)

var _ = Describe("react", func() {

	worldName := "prime"
	worldDescription := "is the main world"

	It("can initialize metadata", func() {
		world := generated.WorldFactory()
		world.Spec().SetName(worldName)

		obj, err := stc.Create(ctx, world)
		Expect(err).To(BeNil())

		newWorld := obj.(generated.World)
		Expect(newWorld).ToNot(BeNil())

		Expect(newWorld.Metadata().Identity()).ToNot(Equal(world.Metadata().Identity()))

		Expect(len(newWorld.Metadata().Created())).ToNot(Equal(0))
		Expect(len(newWorld.Metadata().Updated())).To(Equal(0))

		// log.Println(utils.PP(newWorld))
		time.Sleep(1 * time.Second)
	})

	It("can update metadata", func() {
		world := generated.WorldFactory()
		world.Spec().SetName(worldName)

		obj, err := stc.Update(ctx, generated.WorldIdentity(worldName), world)
		Expect(err).To(BeNil())

		newWorld := obj.(generated.World)
		Expect(len(newWorld.Metadata().Updated())).ToNot(Equal(0))
		Expect(newWorld.Metadata().Identity()).ToNot(
			Equal(world.Metadata().Identity()))
		Expect(newWorld.Metadata().Updated()).ToNot(
			Equal(world.Metadata().Updated()))
		Expect(newWorld.Metadata().Created()).ToNot(
			Equal(world.Metadata().Created()))

		// log.Println(utils.PP(newWorld))
		time.Sleep(1 * time.Second)
	})

	It("can reset status", func() {
		obj, err := stc.Get(ctx, generated.WorldIdentity(worldName))
		Expect(err).To(BeNil())

		world := obj.(generated.World)
		world.Spec().SetName(worldName)
		world.Status().SetDescription(worldDescription)

		// log.Println(utils.PP(world))

		obj, err = stc.Update(ctx, generated.WorldIdentity(worldName), world)
		Expect(err).To(BeNil())

		newWorld := obj.(generated.World)

		// log.Println(utils.PP(world))
		// log.Println(utils.PP(newWorld))

		Expect(newWorld.Metadata().Identity()).To(
			Equal(world.Metadata().Identity()))
		Expect(newWorld.Metadata().Created()).To(
			Equal(world.Metadata().Created()))
		Expect(newWorld.Metadata().Updated()).ToNot(
			Equal(world.Metadata().Updated()))

		Expect(newWorld.Status().Description).ToNot(Equal(
			world.Status().Description()))
	})
})
