package react

import (
	"context"
	"fmt"

	"github.com/wazofski/storz/internal/constants"
	"github.com/wazofski/storz/internal/logger"
	"github.com/wazofski/storz/store"
	"github.com/wazofski/storz/store/options"
)

type reactStore struct {
	Schema           store.SchemaHolder
	Store            store.Store
	Log              logger.Logger
	CallbackRegistry map[string]map[int]Callback
}

const (
	ActionCreate = 1
	ActionUpdate = 2
	ActionDelete = 3
)

type _Register func(d *reactStore) error
type Callback func(store.Object, store.Store) error

func Register(typ string, action int, callback Callback) _Register {
	return func(d *reactStore) error {
		if action < 1 || action > 3 {
			return fmt.Errorf("invalid action %d", action)
		}

		proto := d.Schema.ObjectForKind(typ)
		if proto == nil {
			return fmt.Errorf("invalid type %s", typ)
		}

		_, ok := d.CallbackRegistry[proto.Metadata().Kind()]
		if !ok {
			d.CallbackRegistry[proto.Metadata().Kind()] = make(map[int]Callback)
		}

		_, ok = d.CallbackRegistry[proto.Metadata().Kind()][action]
		if ok {
			return fmt.Errorf("callback for %s %d already set", typ, action)
		}

		d.CallbackRegistry[proto.Metadata().Kind()][action] = callback
		return nil
	}
}

func ReactFactory(data store.Store, callbacks ..._Register) store.Factory {
	return func(schema store.SchemaHolder) (store.Store, error) {
		client := &reactStore{
			Schema:           schema,
			Store:            data,
			Log:              logger.Factory("react"),
			CallbackRegistry: make(map[string]map[int]Callback),
		}

		for _, c := range callbacks {
			err := c(client)
			if err != nil {
				return nil, err
			}
		}

		return client, nil
	}
}

func (d *reactStore) Create(
	ctx context.Context,
	obj store.Object,
	opt ...options.CreateOption) (store.Object, error) {

	if obj == nil {
		return nil, constants.ErrObjectNil
	}

	d.Log.Printf("create %s", obj.PrimaryKey())

	err := d.runCallback(obj, ActionCreate)
	if err != nil {
		return nil, err
	}

	return d.Store.Create(ctx, obj, opt...)
}

func (d *reactStore) Update(
	ctx context.Context,
	identity store.ObjectIdentity,
	obj store.Object,
	opt ...options.UpdateOption) (store.Object, error) {

	if obj == nil {
		return nil, constants.ErrObjectNil
	}

	d.Log.Printf("update %s", identity.Path())
	existing, _ := d.Get(ctx, identity)
	if existing == nil {
		return nil, constants.ErrNoSuchObject
	}

	// update spec
	specHolder := existing.(store.SpecHolder)
	if specHolder != nil {
		specHolder.SpecInternalSet(obj.(store.SpecHolder).SpecInternal())
	}

	err := d.runCallback(existing, ActionUpdate)
	if err != nil {
		return nil, err
	}

	return d.Store.Update(ctx, identity, obj, opt...)
}

func (d *reactStore) Delete(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.DeleteOption) error {

	d.Log.Printf("delete %s", identity.Path())

	existing, _ := d.Get(ctx, identity)
	if existing == nil {
		return constants.ErrNoSuchObject
	}

	err := d.runCallback(existing, ActionDelete)
	if err != nil {
		return err
	}

	return d.Store.Delete(ctx, identity, opt...)
}

func (d *reactStore) Get(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.GetOption) (store.Object, error) {

	d.Log.Printf("get %s", identity.Path())

	return d.Store.Get(ctx, identity, opt...)
}

func (d *reactStore) List(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.ListOption) (store.ObjectList, error) {

	d.Log.Printf("list %s", identity.Type())

	return d.Store.List(ctx, identity, opt...)
}

func (d *reactStore) runCallback(obj store.Object, action int) error {

	_, ok := d.CallbackRegistry[obj.Metadata().Kind()]
	if !ok {
		return nil
	}

	_, ok = d.CallbackRegistry[obj.Metadata().Kind()][action]
	if !ok {
		return nil
	}

	return d.CallbackRegistry[obj.Metadata().Kind()][action](obj, d)
}
