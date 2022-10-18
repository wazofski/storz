package react

import (
	"context"

	"github.com/wazofski/storz/constants"
	"github.com/wazofski/storz/logger"
	"github.com/wazofski/storz/store"
	"github.com/wazofski/storz/store/options"
	"github.com/wazofski/storz/utils"
)

type metaHandlerStore struct {
	Schema store.SchemaHolder
	Store  store.Store
	Log    logger.Logger
}

func MetaHHandlerFactory(data store.Store) store.Factory {
	return func(schema store.SchemaHolder) (store.Store, error) {
		client := &metaHandlerStore{
			Schema: schema,
			Store:  data,
			Log:    logger.Factory("meta handler"),
		}

		return client, nil
	}
}

func (d *metaHandlerStore) Create(
	ctx context.Context,
	obj store.Object,
	opt ...options.CreateOption) (store.Object, error) {

	if obj == nil {
		return nil, constants.ErrObjectNil
	}

	d.Log.Printf("create %s", obj.PrimaryKey())

	ms := obj.Metadata().(store.MetaSetter)

	ms.SetIdentity(store.ObjectIdentityFactory())
	ms.SetCreated(utils.Timestamp())

	return d.Store.Create(ctx, obj, opt...)
}

func (d *metaHandlerStore) Update(
	ctx context.Context,
	identity store.ObjectIdentity,
	obj store.Object,
	opt ...options.UpdateOption) (store.Object, error) {

	if obj == nil {
		return nil, constants.ErrObjectNil
	}

	d.Log.Printf("update %s", identity.Path())

	// update metadata
	ms := obj.Metadata().(store.MetaSetter)
	ms.SetUpdated(utils.Timestamp())

	return d.Store.Update(ctx, identity, obj, opt...)
}

func (d *metaHandlerStore) Delete(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.DeleteOption) error {

	d.Log.Printf("delete %s", identity.Path())

	return d.Store.Delete(ctx, identity, opt...)
}

func (d *metaHandlerStore) Get(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.GetOption) (store.Object, error) {

	d.Log.Printf("get %s", identity.Path())

	return d.Store.Get(ctx, identity, opt...)
}

func (d *metaHandlerStore) List(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.ListOption) (store.ObjectList, error) {

	d.Log.Printf("list %s", identity.Type())

	return d.Store.List(ctx, identity, opt...)
}
