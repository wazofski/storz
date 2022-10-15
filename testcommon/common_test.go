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

		ret, err := clt.Create(ctx, w)

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret.Metadata().Identity())).ToNot(Equal(0))
	})

	It("can GET objects", func() {
		ret, err := clt.Get(ctx, generated.WorldIdentity("abc"))

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

})
