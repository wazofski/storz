package cache

import (
	"context"

	"github.com/wazofski/store"
	"github.com/wazofski/store/memory"
)

type cachedStore struct {
	Schema store.SchemaHolder
	Store  store.Store
	Cache  store.Store
}

func StoreFactory(module string, st store.Store) store.Factory {
	return func(schema store.SchemaHolder) (store.Store, error) {
		client := &cachedStore{
			Schema: schema,
			Store:  st,
			Cache:  store.New(schema, memory.Factory()),
		}

		return client, nil
	}
}

func (d *cachedStore) Create(
	ctx context.Context,
	obj store.Object,
	opt ...store.CreateOption) (store.Object, error) {

	d.Cache.Create(ctx, obj, opt...)

	return d.Store.Create(ctx, obj, opt...)
}

func (d *cachedStore) Update(
	ctx context.Context,
	identity store.ObjectIdentity,
	obj store.Object,
	opt ...store.UpdateOption) (store.Object, error) {

	d.Cache.Update(ctx, identity, obj, opt...)

	return d.Store.Update(ctx, identity, obj, opt...)
}

func (d *cachedStore) Delete(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...store.DeleteOption) error {

	d.Cache.Delete(ctx, identity, opt...)

	return d.Store.Delete(ctx, identity, opt...)
}

func (d *cachedStore) Get(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...store.GetOption) (store.Object, error) {

	ret, _ := d.Cache.Get(ctx, identity, opt...)
	if ret != nil {
		return ret, nil
	}

	return d.Store.Get(ctx, identity, opt...)
}

func (d *cachedStore) List(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...store.ListOption) (store.ObjectList, error) {

	return d.Store.List(ctx, identity, opt...)
}
