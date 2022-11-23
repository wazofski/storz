package cache_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cache", func() {

	ctx := context.Background()

	It("can immediately expire by default", func() {
		// create into cache
		// update inside the real one
		// wait one sec
		// get from cache
		// must match the real one
	})

	It("can override default expire - create", func() {
		// create into cache with 10 min expiration
		// update inside the real one
		// wait one sec
		// get from cache
		// must not match the real one
		// must match the cached one
	})

	It("can override default expire - update", func() {
		// update into cache with 10 min expiration
		// update inside the real one
		// wait one sec
		// get from cache
		// must not match the real one
		// must match the cached one
	})

	It("can expire - get", func() {
		// create into cache
		// update inside the real one
		// wait one sec
		// get from cache
		// must match the real one
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

	It("can list same as real", func() {

	})

})
