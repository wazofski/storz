package react

import (
	"context"
	"fmt"
	"log"
	"time"

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

	if obj == nil {
		return nil, fmt.Errorf("object is nil")
	}

	log.Printf("REACT create %s", obj.PrimaryKey())

	// initialize metadata
	original := d.Schema.ObjectForKind(obj.Metadata().Kind())
	if original == nil {
		return nil, fmt.Errorf("unknown kind %s", obj.Metadata().Kind())
	}

	ms := original.Metadata().(store.MetaSetter)

	ms.SetIdentity(store.ObjectIdentityFactory())
	ms.SetCreated(timestamp())

	// oms := original.(store.MetadataSetter)
	// oms.SetMetadata(ms.(store.Meta))

	// update spec
	specHolder := original.(store.SpecHolder)
	if specHolder != nil {
		specHolder.SpecInternalSet(obj.(store.SpecHolder).SpecInternal())
	}

	return d.Store.Create(ctx, original, opt...)
}

func (d *reactStore) Update(
	ctx context.Context,
	identity store.ObjectIdentity,
	obj store.Object,
	opt ...store.UpdateOption) (store.Object, error) {

	log.Printf("REACT update %s", identity.Path())
	// read the real object
	original, err := d.Store.Get(ctx, identity)

	// if doesn't exist return error
	if err != nil {
		return nil, err
	}
	if original == nil {
		return nil, fmt.Errorf("object %s does not exist", identity)
	}

	// update metadata
	ms := original.Metadata().(store.MetaSetter)
	ms.SetUpdated(timestamp())

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

func timestamp() string {
	return time.Now().Format(time.RFC3339)
}
