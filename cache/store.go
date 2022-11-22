package cache

import (
	"context"
	"time"

	"github.com/wazofski/storz/internal/logger"
	"github.com/wazofski/storz/memory"
	"github.com/wazofski/storz/store"
	"github.com/wazofski/storz/store/options"
)

var log = logger.Factory("rest client")

type cachedStore struct {
	Schema            store.SchemaHolder
	Store             store.Store
	Cache             store.Store
	DefaultExpiration time.Duration
	Policies          map[store.ObjectIdentity]time.Duration
	Modiffies         map[store.ObjectIdentity]time.Time
}

type cacheOptions struct {
	options.CommonOptionHolder
	Expiration time.Duration
}

func newCacheOptions(d *cachedStore) cacheOptions {
	res := cacheOptions{
		CommonOptionHolder: options.CommonOptionHolderFactory(),
		Expiration:         d.DefaultExpiration,
	}

	return res
}

func (d *cacheOptions) CommonOptions() *options.CommonOptionHolder {
	return &d.CommonOptionHolder
}

func Factory(st store.Store, exp ...time.Duration) store.Factory {
	return func(schema store.SchemaHolder) (store.Store, error) {
		client := &cachedStore{
			Schema:    schema,
			Store:     st,
			Cache:     store.New(schema, memory.Factory()),
			Policies:  make(map[store.ObjectIdentity]time.Duration),
			Modiffies: make(map[store.ObjectIdentity]time.Time),
		}

		if len(exp) > 0 {
			client.DefaultExpiration = exp[0]
		}

		return client, nil
	}
}

func (d *cachedStore) Create(
	ctx context.Context,
	obj store.Object,
	opt ...options.CreateOption) (store.Object, error) {

	copt := newCacheOptions(d)
	for _, o := range opt {
		err := o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	ret, err := d.Cache.Create(ctx, obj, opt...)
	if copt.Expiration > 0 && err == nil {
		d.Policies[ret.Metadata().Identity()] = copt.Expiration
		d.Modiffies[ret.Metadata().Identity()] = time.Now()
	}

	return d.Store.Create(ctx, obj, opt...)
}

func (d *cachedStore) Update(
	ctx context.Context,
	identity store.ObjectIdentity,
	obj store.Object,
	opt ...options.UpdateOption) (store.Object, error) {

	copt := newCacheOptions(d)
	for _, o := range opt {
		err := o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	ret, err := d.Cache.Update(ctx, identity, obj, opt...)
	if copt.Expiration > 0 && err == nil {
		d.Policies[ret.Metadata().Identity()] = copt.Expiration
		d.Modiffies[ret.Metadata().Identity()] = time.Now()
	}

	return d.Store.Update(ctx, identity, obj, opt...)
}

func (d *cachedStore) Delete(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.DeleteOption) error {

	err := d.Cache.Delete(ctx, identity, opt...)
	if err == nil {
		delete(d.Policies, identity)
		delete(d.Modiffies, identity)
	}

	return d.Store.Delete(ctx, identity, opt...)
}

func (d *cachedStore) Get(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.GetOption) (store.Object, error) {

	has_expired := false
	exp := time.Duration(0)
	exp, ok := d.Policies[identity]
	if !ok {
		exp = d.DefaultExpiration
	}

	if exp > 0 {
		expt := d.Modiffies[identity].Add(exp)
		if expt.Before(time.Now()) {
			has_expired = true
		}
	}

	cached, cached_err := d.Cache.Get(ctx, identity)
	if has_expired || cached == nil {
		ret, err := d.Store.Get(ctx, identity, opt...)
		if ret != nil && err == nil {
			d.Policies[identity] = exp
			d.Modiffies[identity] = time.Now()
			if cached == nil {
				d.Cache.Create(ctx, ret)
			} else {
				d.Cache.Update(ctx, identity, ret)
			}
		}
		return ret, err
	}

	return cached, cached_err
}

func (d *cachedStore) List(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.ListOption) (store.ObjectList, error) {

	return d.Store.List(ctx, identity, opt...)
}
