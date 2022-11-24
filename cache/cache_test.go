package cache_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/storz/cache"
	"github.com/wazofski/storz/generated"
)

var _ = Describe("cache", func() {

	ctx := context.Background()

	worldName := "abc"
	worldDesc := "def"

	It("can immediately expire by default", func() {
		// create into cache
		world := generated.WorldFactory()
		world.Spec().SetName(worldName)
		ret, err := cached.Create(ctx, world)

		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())

		// update inside the real one
		world = ret.(generated.World)
		world.Spec().SetDescription(worldDesc)
		ret, err = mainst.Update(ctx, world.Metadata().Identity(), world)
		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())

		// get from cache
		ret, err = cached.Get(ctx, world.Metadata().Identity())
		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())
		world = ret.(generated.World)

		// must match the real one
		Expect(world.Spec().Description()).To(Equal(""))

		// wait one sec
		time.Sleep(1 * time.Second)

		// cached one must be expired to the real one will be fetched
		// get from cache
		ret, err = cached.Get(ctx, world.Metadata().Identity())
		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())
		world = ret.(generated.World)

		// must match the real one
		Expect(world.Spec().Description()).To(Equal(worldDesc))

		err = cached.Delete(ctx, world.Metadata().Identity())
		Expect(err).To(BeNil())
	})

	It("can override default expire - create", func() {
		// create into cache with 10 min expiration
		world := generated.WorldFactory()
		world.Spec().SetName(worldName)
		ret, err := cached.Create(ctx, world, cache.Expire(10*time.Minute))

		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())

		// update inside the real one
		world = ret.(generated.World)
		world.Spec().SetDescription(worldDesc)
		ret, err = mainst.Update(ctx, world.Metadata().Identity(), world)
		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())

		// get from cache
		ret, err = cached.Get(ctx, world.Metadata().Identity())
		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())
		world = ret.(generated.World)

		// must match the cached one
		Expect(world.Spec().Description()).To(Equal(""))

		// wait one sec
		time.Sleep(1 * time.Second)

		// cached one must not be expired so the cached one will be fetched
		// get from cache
		ret, err = cached.Get(ctx, world.Metadata().Identity())
		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())
		world = ret.(generated.World)

		// must match the cached one
		Expect(world.Spec().Description()).To(Equal(""))
	})

	It("can override default expire - update", func() {
		// update into cache with 10 min expiration
		ret, err := cached.Get(ctx,
			generated.WorldIdentity(worldName))

		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())

		world := ret.(generated.World)
		world.Spec().SetDescription(worldDesc)

		ret, err = cached.Update(ctx,
			world.Metadata().Identity(),
			world,
			cache.Expire(10*time.Minute))

		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())

		// update inside the real one
		world = ret.(generated.World)
		world.Spec().SetDescription("")

		ret, err = mainst.Update(ctx, world.Metadata().Identity(), world)
		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())

		// get from cache
		ret, err = cached.Get(ctx, world.Metadata().Identity())
		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())
		world = ret.(generated.World)

		// must match the cached one
		Expect(world.Spec().Description()).To(Equal(worldDesc))

		// wait one sec
		time.Sleep(1 * time.Second)

		// cached one must not be expired so the cached one will be fetched
		// get from cache
		ret, err = cached.Get(ctx, world.Metadata().Identity())
		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())
		world = ret.(generated.World)

		// must match the cached one
		Expect(world.Spec().Description()).To(Equal(worldDesc))
	})

	It("can expire - create", func() {
		// create into cache with 2 sec expiration
		// update inside the real one
		// wait one sec
		// get from cache
		// must match the cached one

		// wait one sec
		// get from cache
		// must match the real one
	})

	It("can expire - update", func() {
		// update into cache with 2 sec expiration
		// update inside the real one
		// wait one sec
		// get from cache
		// must match the cached one

		// wait one sec
		// get from cache
		// must match the real one
	})

	It("can delete unexpired", func() {
		// update into cache with 2 sec expiration
		// update inside the real one
		// wait one sec
		// get from cache
		// must match the cached one

		// wait one sec
		// get from cache
		// must match the real one
	})

	It("can delete expired", func() {
		// update into cache with 2 sec expiration
		// update inside the real one
		// wait one sec
		// get from cache
		// must match the cached one

		// wait one sec
		// get from cache
		// must match the real one
	})

})
