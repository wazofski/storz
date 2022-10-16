package react

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wazofski/store"
)

type metaHandlerStore struct {
	Schema store.SchemaHolder
	Store  store.Store
}

func MetaHHandlerFactory(data store.Store) store.Factory {
	return func(schema store.SchemaHolder) (store.Store, error) {
		client := &metaHandlerStore{
			Schema: schema,
			Store:  data,
		}

		return client, nil
	}
}

func (d *metaHandlerStore) Create(
	ctx context.Context,
	obj store.Object,
	opt ...store.CreateOption) (store.Object, error) {

	if obj == nil {
		return nil, fmt.Errorf("object is nil")
	}

	log.Printf("META HANDLER create %s", obj.PrimaryKey())

	ms := obj.Metadata().(store.MetaSetter)

	ms.SetIdentity(store.ObjectIdentityFactory())
	ms.SetCreated(timestamp())

	return d.Store.Create(ctx, obj, opt...)
}

func (d *metaHandlerStore) Update(
	ctx context.Context,
	identity store.ObjectIdentity,
	obj store.Object,
	opt ...store.UpdateOption) (store.Object, error) {

	if obj == nil {
		return nil, fmt.Errorf("object is nil")
	}

	log.Printf("META HANDLER update %s", identity.Path())

	// update metadata
	ms := obj.Metadata().(store.MetaSetter)
	ms.SetUpdated(timestamp())

	return d.Store.Update(ctx, identity, obj, opt...)
}

func (d *metaHandlerStore) Delete(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...store.DeleteOption) error {

	log.Printf("META HANDLER delete %s", identity.Path())

	return d.Store.Delete(ctx, identity, opt...)
}

func (d *metaHandlerStore) Get(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...store.GetOption) (store.Object, error) {

	log.Printf("META HANDLER get %s", identity.Path())

	return d.Store.Get(ctx, identity, opt...)
}

func (d *metaHandlerStore) List(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...store.ListOption) (store.ObjectList, error) {

	log.Printf("META HANDLER list %s", identity.Type())

	return d.Store.List(ctx, identity, opt...)
}

func timestamp() string {
	return time.Now().Format(time.RFC3339)
}
