package memory

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/wazofski/store"
)

type reactStore struct {
	Schema store.SchemaHolder
	Store  store.Store
}

func Factory(data store.Store) store.Factory {
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

	log.Printf("REACT create %s", obj.PrimaryKey())

	// initialize metadata
	obj.Metadata().SetIdentity(store.ObjectIdentity(uuid.New().String()))

	// reset status to nothing

	return d.Store.Create(ctx, obj, opt...)
}

func (d *reactStore) Update(
	ctx context.Context,
	identity store.ObjectIdentity,
	obj store.Object,
	opt ...store.UpdateOption) (store.Object, error) {

	log.Printf("REACT update %s", identity.Path())
	// read the real object

	// if doesn't exist return error

	// reset and update metadata

	// reset status

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
