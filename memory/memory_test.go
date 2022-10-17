package memory_test

import (
	"sort"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/store"
	"github.com/wazofski/store/generated"
	"github.com/wazofski/store/options"
)

var _ = Describe("memory", func() {
	worldName := "c137"
	anotherWorldName := "j19zeta7"
	worldDescription := "the world of argo"
	newWorldDescription := "is only beginning"
	worldId := store.ObjectIdentity("")

	It("cannot GET nonexistent objects", func() {
		ret, err := client.Get(
			ctx,
			generated.WorldIdentity(anotherWorldName),
		)

		Expect(ret).To(BeNil())
		Expect(err).ToNot(BeNil())
	})

	It("can LIST empty lists", func() {
		ret, err := client.List(
			ctx, generated.WorldIdentity(""))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret)).To(Equal(0))
	})

	It("can POST objects", func() {
		w := generated.WorldFactory()

		w.Spec().SetName(worldName)
		w.Spec().SetDescription(worldDescription)

		ret, err := client.Create(ctx, w)
		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret.Metadata().Identity())).ToNot(Equal(0))
	})

	It("cannot LIST and FILTER BY nonexistent props", func() {
		ret, err := client.List(
			ctx, generated.WorldIdentity(""),
			options.PropFilter("metadata.askdjhasd", "asdsadas"))

		Expect(err).ToNot(BeNil())
		Expect(ret).To(BeNil())
	})

	It("can LIST single object", func() {
		ret, err := client.List(
			ctx, generated.WorldIdentity(""))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret)).To(Equal(1))

		world := ret[0].(generated.World)
		Expect(world.Spec().Name()).To(Equal(worldName))
		Expect(world.Spec().Description()).To(Equal(worldDescription))
	})

	It("cannot LIST specific object", func() {
		// world exists
		ret, err := client.List(
			ctx, generated.WorldIdentity(worldName))

		Expect(ret).To(BeNil())
		Expect(err).ToNot(BeNil())
	})

	It("cannot LIST specific nonexistent object", func() {
		// another world does not exist
		ret, err := client.List(
			ctx, generated.WorldIdentity(anotherWorldName))

		Expect(ret).To(BeNil())
		Expect(err).ToNot(BeNil())
	})

	It("can POST other objects", func() {
		w := generated.SecondWorldFactory()

		w.Spec().SetName(anotherWorldName)
		w.Spec().SetDescription(newWorldDescription)

		ret, err := client.Create(ctx, w)
		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret.Metadata().Identity())).ToNot(Equal(0))

		ret, err = client.Get(ctx,
			generated.SecondWorldIdentity(anotherWorldName))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret.Metadata().Identity())).ToNot(Equal(0))

		w = generated.WorldFactory()

		w.Spec().SetName(anotherWorldName)
		w.Spec().SetDescription(newWorldDescription)

		ret, err = client.Create(ctx, w)
		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret.Metadata().Identity())).ToNot(Equal(0))
	})

	It("can LIST multiple objects", func() {
		ret, err := client.List(
			ctx, generated.WorldIdentity(""))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret)).To(Equal(2))

		sort.Slice(ret, func(i, j int) bool {
			return ret[i].(generated.World).Spec().Name() < ret[j].(generated.World).Spec().Name()
		})

		world := ret[0].(generated.World)
		Expect(world.Spec().Name()).To(Equal(worldName))
		Expect(world.Spec().Description()).To(Equal(worldDescription))

		world2 := ret[1].(generated.World)
		Expect(world2.Spec().Name()).To(Equal(anotherWorldName))
		Expect(world2.Spec().Description()).To(Equal(newWorldDescription))
	})

	It("can LIST and sort multiple objects", func() {
		ret, err := client.List(
			ctx, generated.WorldIdentity(""),
			options.OrderBy("spec.name"))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret)).To(Equal(2))

		world := ret[0].(generated.World)
		Expect(world.Spec().Name()).To(Equal(worldName))
		Expect(world.Spec().Description()).To(Equal(worldDescription))

		world2 := ret[1].(generated.World)
		Expect(world2.Spec().Name()).To(Equal(anotherWorldName))
		Expect(world2.Spec().Description()).To(Equal(newWorldDescription))

		ret, err = client.List(
			ctx, generated.WorldIdentity(""),
			options.OrderBy("spec.name"),
			options.OrderDescending())

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret)).To(Equal(2))

		world = ret[1].(generated.World)
		world2 = ret[0].(generated.World)
		Expect(world.Spec().Name()).To(Equal(worldName))
		Expect(world2.Spec().Name()).To(Equal(anotherWorldName))
	})

	It("can LIST and paginate multiple objects", func() {
		w := generated.SecondWorldFactory()
		w.Spec().SetName(worldName)
		w.Spec().SetDescription(worldDescription)
		client.Create(ctx, w)

		ret, err := client.List(
			ctx, generated.SecondWorldIdentity(""),
			options.OrderBy("spec.name"),
			options.PageSize(1))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret)).To(Equal(1))

		world := ret[0].(generated.SecondWorld)
		Expect(world.Spec().Name()).To(Equal(worldName))
		Expect(world.Spec().Description()).To(Equal(worldDescription))

		ret, err = client.List(
			ctx, generated.SecondWorldIdentity(""),
			options.OrderBy("spec.name"),
			options.PageSize(1),
			options.PageOffset(1))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret)).To(Equal(1))

		world = ret[0].(generated.SecondWorld)
		Expect(world.Spec().Name()).To(Equal(anotherWorldName))
	})

	It("can LIST and FILTER", func() {
		ret, err := client.List(
			ctx, generated.SecondWorldIdentity(""),
			options.PropFilter("spec.name", worldName))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret)).To(Equal(1))

		world := ret[0].(generated.SecondWorld)
		Expect(world.Spec().Name()).To(Equal(worldName))
		Expect(world.Spec().Description()).To(Equal(worldDescription))
		worldId = world.Metadata().Identity()
	})

	It("can LIST and FILTER BY ID", func() {
		ret, err := client.List(
			ctx, generated.SecondWorldIdentity(""),
			options.PropFilter("metadata.identity", string(worldId)))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret)).To(Equal(1))

		world := ret[0].(generated.SecondWorld)
		Expect(world.Spec().Name()).To(Equal(worldName))
		Expect(world.Spec().Description()).To(Equal(worldDescription))
	})

	It("cannot double POST objects", func() {
		w := generated.SecondWorldFactory()

		w.Spec().SetName(anotherWorldName)
		w.Spec().SetDescription(newWorldDescription)

		ret, err := client.Create(ctx, w)
		Expect(err).ToNot(BeNil())
		Expect(ret).To(BeNil())
	})

	It("can GET objects", func() {
		ret, err := client.Get(
			ctx, generated.SecondWorldIdentity(anotherWorldName))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		world := ret.(generated.SecondWorld)
		Expect(world.Spec().Name()).To(Equal(anotherWorldName))
		Expect(world.Spec().Description()).To(Equal(newWorldDescription))

		worldId = world.Metadata().Identity()
	})

	It("can GET objects by ID", func() {
		ret, err := client.Get(ctx, worldId)

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		world := ret.(generated.SecondWorld)
		Expect(world.Metadata().Identity()).To(Equal(worldId))
		Expect(world.Spec().Name()).To(Equal(anotherWorldName))
		Expect(world.Spec().Description()).To(Equal(newWorldDescription))
	})

	It("can PUT objects", func() {
		ret, err := client.Get(
			ctx,
			generated.SecondWorldIdentity(anotherWorldName))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		world := ret.(generated.SecondWorld)
		Expect(world).ToNot(BeNil())

		world.Spec().SetDescription(newWorldDescription)
		ret, err = client.Update(
			ctx,
			generated.SecondWorldIdentity(anotherWorldName),
			world)

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		ret, err = client.Get(
			ctx,
			generated.SecondWorldIdentity(anotherWorldName))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		world = ret.(generated.SecondWorld)
		Expect(world).ToNot(BeNil())
		Expect(world.Spec().Description()).To(Equal(newWorldDescription))
	})

	It("can DELETE objects", func() {
		err := client.Delete(
			ctx, generated.SecondWorldIdentity(anotherWorldName))

		Expect(err).To(BeNil())
	})

	It("cannot DELETE nonexistent objects", func() {
		err := client.Delete(
			ctx, generated.SecondWorldIdentity(anotherWorldName))

		Expect(err).ToNot(BeNil())
	})

	It("cannot PUT nonexistent objects", func() {
		w := generated.SecondWorldFactory()

		w.Spec().SetName(anotherWorldName)
		w.Spec().SetDescription(worldDescription)

		ret, err := client.Update(
			ctx,
			generated.SecondWorldIdentity(anotherWorldName),
			w)

		Expect(ret).To(BeNil())
		Expect(err).ToNot(BeNil())
	})
})
