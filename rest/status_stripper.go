package rest

import (
	"context"
	"fmt"

	"github.com/wazofski/storz/internal/constants"
	"github.com/wazofski/storz/internal/logger"
	"github.com/wazofski/storz/store"
	"github.com/wazofski/storz/store/options"
)

type statusStripperStore struct {
	Schema store.SchemaHolder
	Store  store.Store
	Log    logger.Logger
}

func _StatusStripperFactory(data store.Store) store.Factory {
	return func(schema store.SchemaHolder) (store.Store, error) {
		client := &statusStripperStore{
			Schema: schema,
			Store:  data,
			Log:    logger.Factory("status stripper"),
		}

		return client, nil
	}
}

func (d *statusStripperStore) Create(
	ctx context.Context,
	obj store.Object,
	opt ...options.CreateOption) (store.Object, error) {

	if obj == nil {
		return nil, constants.ErrObjectNil
	}

	d.Log.Printf("create %s", obj.PrimaryKey())

	// initialize metadata
	original := d.Schema.ObjectForKind(obj.Metadata().Kind())
	if original == nil {
		return nil, fmt.Errorf("unknown kind %s", obj.Metadata().Kind())
	}

	// update spec
	specHolder := original.(store.SpecHolder)
	if specHolder != nil {
		specHolder.SpecInternalSet(obj.(store.SpecHolder).SpecInternal())
	}

	return d.Store.Create(ctx, original, opt...)
}

func (d *statusStripperStore) Update(
	ctx context.Context,
	identity store.ObjectIdentity,
	obj store.Object,
	opt ...options.UpdateOption) (store.Object, error) {

	if obj == nil {
		return nil, constants.ErrObjectNil
	}

	d.Log.Printf("update %s", identity.Path())
	// read the real object
	original, err := d.Store.Get(ctx, identity)

	// if doesn't exist return error
	if err != nil {
		return nil, err
	}
	if original == nil {
		return nil, constants.ErrNoSuchObject
	}

	// update spec
	specHolder := original.(store.SpecHolder)
	if specHolder != nil && obj != nil {
		objSpecHolder := obj.(store.SpecHolder)
		if objSpecHolder != nil {
			specHolder.SpecInternalSet(
				objSpecHolder.SpecInternal())
		}
	}

	return d.Store.Update(ctx, identity, original, opt...)
}

func (d *statusStripperStore) Delete(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.DeleteOption) error {

	d.Log.Printf("delete %s", identity.Path())

	return d.Store.Delete(ctx, identity, opt...)
}

func (d *statusStripperStore) Get(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.GetOption) (store.Object, error) {

	d.Log.Printf("get %s", identity.Path())

	return d.Store.Get(ctx, identity, opt...)
}

func (d *statusStripperStore) List(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.ListOption) (store.ObjectList, error) {

	d.Log.Printf("list %s", identity.Type())

	return d.Store.List(ctx, identity, opt...)
}
