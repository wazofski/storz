package client_test

import (
	"log"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/storz/client"
	"github.com/wazofski/storz/generated"
	"github.com/wazofski/storz/store"
	"github.com/wazofski/storz/store/options"
)

var _ = Describe("client", func() {
	worldName := "c137"

	It("can specify client.Headers", func() {
		_, err := stc.Get(
			ctx, "",
			client.Header("setting", "a client.Header"),
			client.Header("setting", "another client.Header"),
		)
		Expect(err).ToNot(BeNil())
		log.Printf("expected error: %s", err)

		_, err = stc.List(
			ctx, "",
			options.PropFilter("metadata.ID", "value"),
			client.Header("setting", "a client.Header"),
			client.Header("another setting", "another client.Header"),
		)
		Expect(err).ToNot(BeNil())
		log.Printf("expected error: %s", err)

		_, err = stc.Update(
			ctx, "", nil,
			// Options for other APIs are not accepted
			// options.PropFilter("metadata.ID", "value"),
			client.Header("setting", "another client.Header"),
		)
		Expect(err).ToNot(BeNil())
		log.Printf("expected error: %s", err)
	})

	It("cannot GET non-allowed", func() {
		ret, err := stc.Get(
			ctx, generated.ThirdWorldIdentity(worldName))

		Expect(ret).To(BeNil())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("http 405"))

		ret, err = stc.Get(
			ctx,
			store.ObjectIdentity("id/aliksjdlsakjdaslkjdaslkj"))

		Expect(ret).To(BeNil())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("http 405"))
	})

	It("cannot CREATE non-allowed", func() {
		w := generated.ThirdWorldFactory()
		ret, err := stc.Create(ctx, w)

		Expect(ret).To(BeNil())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("http 405"))
	})

	It("cannot UPDATE non-allowed", func() {
		w := generated.ThirdWorldFactory()
		ret, err := stc.Update(ctx,
			generated.ThirdWorldIdentity(worldName), w)

		Expect(ret).To(BeNil())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("http 405"))

		ret, err = stc.Update(ctx,
			store.ObjectIdentity("id/aliksjdlsakjdaslkjdaslkj"), w)

		Expect(ret).To(BeNil())
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("http 405"))
	})

	It("cannot DELETE non-allowed", func() {
		err := stc.Delete(
			ctx, generated.ThirdWorldIdentity(worldName))

		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("http 405"))

		err = stc.Delete(
			ctx,
			store.ObjectIdentity("id/aliksjdlsakjdaslkjdaslkj"))

		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("http 405"))
	})

	It("cannot LIST non-allowed", func() {
		ret, err := stc.List(
			ctx, generated.ThirdWorldIdentity(""))

		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("http 405"))

		Expect(len(ret)).To(Equal(0))
	})

})
