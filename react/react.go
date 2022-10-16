package react

import (
	"context"
	"fmt"
	"log"

	"github.com/wazofski/store"
)

type reactStore struct {
	Schema store.SchemaHolder
	Store  store.Store
}

func ReactFactory(data store.Store) store.Factory {
	return func(schema store.SchemaHolder) (store.Store, error) {
		client := &reactStore{
			Schema: schema,
			Store:  data,
		}

		return client, nil
	}
}

func (d *reactStore) Create(
	ctx context.Context,
	obj store.Object,
	opt ...store.CreateOption) (store.Object, error) {

	if obj == nil {
		return nil, fmt.Errorf("object is nil")
	}

	log.Printf("REACT create %s", obj.PrimaryKey())

	return d.Store.Create(ctx, obj, opt...)
}

func (d *reactStore) Update(
	ctx context.Context,
	identity store.ObjectIdentity,
	obj store.Object,
	opt ...store.UpdateOption) (store.Object, error) {

	if obj == nil {
		return nil, fmt.Errorf("object is nil")
	}

	log.Printf("REACT update %s", identity.Path())

	return d.Store.Update(ctx, identity, obj, opt...)
}

func (d *reactStore) Delete(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...store.DeleteOption) error {

	log.Printf("REACT delete %s", identity.Path())

	return d.Store.Delete(ctx, identity, opt...)
}

func (d *reactStore) Get(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...store.GetOption) (store.Object, error) {

	log.Printf("REACT get %s", identity.Path())

	return d.Store.Get(ctx, identity, opt...)
}

func (d *reactStore) List(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...store.ListOption) (store.ObjectList, error) {

	log.Printf("REACT list %s", identity.Type())

	return d.Store.List(ctx, identity, opt...)
}
