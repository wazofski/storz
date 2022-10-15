package common_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/store"
	"github.com/wazofski/store/generated"
)

var _ = Describe("common", func() {

	It("can POST objects", func() {
		w := generated.WorldFactory()

		w.Spec().SetName("abc")

		ret, err := clt.Create(ctx, w)

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret.Metadata().Identity())).ToNot(Equal(0))
	})

	It("can GET objects", func() {
		ret, err := clt.Get(ctx,
			generated.WorldIdentity("abc"))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret.Metadata().Identity())).ToNot(Equal(0))

		world := ret.(generated.World)
		Expect(world).ToNot(BeNil())
	})

	It("cannot double POST objects", func() {
		w := generated.WorldFactory()

		w.Spec().SetName("abc")

		ret, err := clt.Create(ctx, w)

		Expect(err).ToNot(BeNil())
		Expect(ret).To(BeNil())
	})

	It("can PUT objects", func() {
		w := generated.WorldFactory()

		w.Spec().SetName("abc")
		w.Spec().SetDescription("def")

		ret, err := clt.Update(ctx,
			generated.WorldIdentity("abc"),
			w)

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		world := ret.(generated.World)
		Expect(world).ToNot(BeNil())
		Expect(world.Spec().Description()).To(Equal("def"))
	})

	It("can PUT change naming props", func() {
		w := generated.WorldFactory()

		w.Spec().SetName("def")

		ret, err := clt.Update(ctx,
			generated.WorldIdentity("abc"),
			w)

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		world := ret.(generated.World)
		Expect(world).ToNot(BeNil())
		Expect(world.Spec().Name()).To(Equal("def"))

		ret, err = clt.Get(ctx,
			generated.WorldIdentity("abc"))

		Expect(err).ToNot(BeNil())
		Expect(ret).To(BeNil())
	})

	It("can PUT objects BY ID", func() {
		ret, err := clt.Get(ctx,
			generated.WorldIdentity("def"))
		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		world := ret.(generated.World)
		Expect(world).ToNot(BeNil())
		world.Spec().SetDescription("qqq")

		ret, err = clt.Update(ctx,
			world.Metadata().Identity(), world)
		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		world = ret.(generated.World)
		Expect(world).ToNot(BeNil())
		Expect(world.Spec().Description()).To(Equal("qqq"))
	})

	It("cannot PUT non-existent objects", func() {
		world := generated.WorldFactory()
		Expect(world).ToNot(BeNil())
		world.Spec().SetName("zxcxzcxz")

		ret, err := clt.Update(ctx,
			generated.WorldIdentity("zcxzcxzc"), world)
		Expect(err).ToNot(BeNil())
		Expect(ret).To(BeNil())
	})

	It("cannot PUT non-existent objects BY ID", func() {
		world := generated.WorldFactory()
		world.Spec().SetName("zxcxzcxz")

		ret, err := clt.Update(ctx,
			world.Metadata().Identity(), world)
		Expect(err).ToNot(BeNil())
		Expect(ret).To(BeNil())
	})

	It("cannot PUT objects of wrong type", func() {
		world := generated.SecondWorldFactory()
		world.Spec().SetName("zxcxzcxz")

		ret, err := clt.Update(ctx,
			generated.WorldIdentity("qwe"), world)
		Expect(err).ToNot(BeNil())
		Expect(ret).To(BeNil())
	})

	It("can GET objects", func() {
		ret, err := clt.Get(ctx,
			generated.WorldIdentity("def"))
		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		world := ret.(generated.World)
		Expect(world).ToNot(BeNil())
	})

	It("can GET objects BY ID", func() {
		ret, err := clt.Get(ctx,
			generated.WorldIdentity("def"))
		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		world := ret.(generated.World)
		Expect(world).ToNot(BeNil())

		ret, err = clt.Get(ctx,
			world.Metadata().Identity())
		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		world = ret.(generated.World)
		Expect(world).ToNot(BeNil())
	})

	It("cannot GET non-existent objects", func() {
		ret, err := clt.Get(ctx,
			generated.WorldIdentity("zxcxzczx"))
		Expect(err).ToNot(BeNil())
		Expect(ret).To(BeNil())
	})

	It("cannot GET non-existent objects BY ID", func() {
		ret, err := clt.Get(ctx,
			store.ObjectIdentity("id/kjjakjjsadldkjalkdajs"))
		Expect(err).ToNot(BeNil())
		Expect(ret).To(BeNil())
	})

	It("can DELETE objects", func() {
		w := generated.WorldFactory()
		w.Spec().SetName("tobedeleted")

		ret, err := clt.Create(ctx, w)
		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		err = clt.Delete(ctx,
			generated.WorldIdentity(w.Spec().Name()))
		Expect(err).To(BeNil())

		_, err = clt.Get(ctx,
			generated.WorldIdentity(w.Spec().Name()))
		Expect(err).ToNot(BeNil())
	})

	It("can DELETE objects BT ID", func() {
		w := generated.WorldFactory()
		w.Spec().SetName("tobedeleted")

		ret, err := clt.Create(ctx, w)
		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		w = ret.(generated.World)

		err = clt.Delete(ctx, w.Metadata().Identity())
		Expect(err).To(BeNil())

		_, err = clt.Get(ctx, w.Metadata().Identity())
		Expect(err).ToNot(BeNil())
	})

	It("cannot DELETE non-existent objects", func() {
		err := clt.Delete(ctx,
			generated.WorldIdentity("akjsdhsajkhdaskjh"))
		Expect(err).ToNot(BeNil())
	})

	It("cannot DELETE non-existent objects BY ID", func() {
		err := clt.Delete(ctx,
			store.ObjectIdentity("id/kjjakjjsadldkjalkdajs"))
		Expect(err).ToNot(BeNil())
	})

	It("cannot GET nil identity", func() {
		_, err := clt.Get(ctx, "")
		Expect(err).ToNot(BeNil())
	})

	It("cannot CREATE nil object", func() {
		_, err := clt.Create(ctx, nil)
		Expect(err).ToNot(BeNil())
	})

	It("cannot PUT nil identity", func() {
		_, err := clt.Update(ctx,
			"", generated.WorldFactory())
		Expect(err).ToNot(BeNil())
	})

	It("cannot PUT nil object", func() {
		_, err := clt.Update(ctx,
			generated.WorldIdentity("qwe"), nil)
		Expect(err).ToNot(BeNil())
	})

	It("cannot DELETE nil identity", func() {
		err := clt.Delete(ctx, "")
		Expect(err).ToNot(BeNil())
	})

})
