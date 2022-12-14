package react_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/storz/generated"
	"github.com/wazofski/storz/store"
)

func WorldCreateCb(obj store.Object, str store.Store) error {
	world := obj.(generated.World)
	world.Status().SetDescription("abc")

	return nil
}

func WorldUpdateCb(obj store.Object, str store.Store) error {
	anotherWorld := generated.SecondWorldFactory()
	anotherWorld.Spec().SetName("def")

	_, err := str.Create(context.Background(), anotherWorld)
	return err
}

func WorldDeleteCb(obj store.Object, str store.Store) error {
	return fmt.Errorf("cannot delete")
}

var _ = Describe("react", func() {

	It("can set status on CREATE", func() {
		world := generated.WorldFactory()
		world.Spec().SetName("abc")

		ret, err := str.Create(ctx, world)

		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())

		world = ret.(generated.World)

		ret, err = str.Get(ctx, world.Metadata().Identity())
		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())

		world = ret.(generated.World)
		Expect(world.Status().Description()).To(Equal("abc"))
	})

	It("can creat objects on UPDATE", func() {
		ret, err := str.Get(ctx, generated.WorldIdentity("abc"))
		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())

		world := ret.(generated.World)
		world.Spec().SetDescription("qwe")
		ret, err = str.Update(ctx, generated.WorldIdentity("abc"), world)
		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())

		ret, err = str.Get(ctx, generated.SecondWorldIdentity("def"))
		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())
	})

	It("can reject DELETE", func() {
		err := str.Delete(ctx, generated.WorldIdentity("abc"))
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("cannot delete"))

		ret, err := str.Get(ctx, generated.WorldIdentity("abc"))
		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())
	})
})
