package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/wazofski/store"
)

type memoryStore struct {
	Schema        store.SchemaHolder
	IdentityIndex map[string]*store.Object
	PrimaryIndex  map[string]map[string]*store.Object
}

func Factory() store.Factory {
	return func(schema store.SchemaHolder) (store.Store, error) {
		client := memoryStore{
			Schema:        schema,
			IdentityIndex: make(map[string]*store.Object),
			PrimaryIndex:  make(map[string]map[string]*store.Object),
		}

		return client, nil
	}
}

func (d memoryStore) Create(
	ctx context.Context,
	obj store.Object,
	opt ...store.CreateOption) (store.Object, error) {

	var err error
	copt := store.CommonOptionHolder{}
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	if obj == nil {
		return nil, fmt.Errorf("object is nil")
	}

	lk := strings.ToLower(obj.Metadata().Kind())
	path := fmt.Sprintf("%s/%s", lk, obj.PrimaryKey())
	existing, _ := d.Get(ctx, store.ObjectIdentity(path))

	if existing != nil {
		return nil, fmt.Errorf("object already exists")
	}

	clone := obj.Clone()

	log.Printf("creating %s", obj.Metadata().Identity())
	log.Printf("path %s", obj.Metadata().Identity().Path())

	d.IdentityIndex[obj.Metadata().Identity().Path()] = &clone
	if d.PrimaryIndex[lk] == nil {
		d.PrimaryIndex[lk] = make(map[string]*store.Object)
	}

	d.PrimaryIndex[lk][obj.PrimaryKey()] = &clone

	return clone, nil
}

func (d memoryStore) Update(
	ctx context.Context,
	identity store.ObjectIdentity,
	obj store.Object,
	opt ...store.UpdateOption) (store.Object, error) {

	log.Printf("updating %s", identity)

	var err error
	copt := store.CommonOptionHolder{}
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	if obj == nil {
		return nil, fmt.Errorf("object is nil")
	}

	existing, _ := d.Get(ctx, identity)
	if existing == nil {
		return nil, fmt.Errorf("object %s does not exist", identity)
	}

	clone := obj.Clone()

	d.IdentityIndex[obj.Metadata().Identity().Path()] = &clone
	lk := strings.ToLower(existing.Metadata().Kind())
	d.PrimaryIndex[lk][existing.PrimaryKey()] = nil

	lk = strings.ToLower(obj.Metadata().Kind())
	d.PrimaryIndex[lk][obj.PrimaryKey()] = &clone

	return clone, err
}

func (d memoryStore) Delete(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...store.DeleteOption) error {

	log.Printf("deleting %s", identity)

	var err error
	copt := store.CommonOptionHolder{}
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return err
		}
	}

	existing, _ := d.Get(ctx, identity)
	if existing == nil {
		return fmt.Errorf("object %s does not exist", identity)
	}

	d.IdentityIndex[identity.Path()] = nil
	lk := strings.ToLower(existing.Metadata().Kind())
	d.PrimaryIndex[lk][existing.PrimaryKey()] = nil

	return nil
}

func (d memoryStore) Get(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...store.GetOption) (store.Object, error) {

	log.Printf("getting %s", identity)

	var err error
	copt := store.CommonOptionHolder{}
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("...GET identity index size: %d", len(d.IdentityIndex))

	ret := d.IdentityIndex[identity.Path()]
	if ret != nil {
		return *ret, nil
	}

	tokens := strings.Split(identity.Path(), "/")
	if len(tokens) == 2 {
		lk := strings.ToLower(tokens[0])
		km := d.PrimaryIndex[lk]
		if km != nil {
			log.Printf("...GET type index exists with %d records", len(km))
			ret = km[tokens[1]]
			if ret != nil {
				return *ret, nil
			}
		}
	}

	return nil, fmt.Errorf("object %s does not exist", identity)
}

func (d memoryStore) List(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...store.ListOption) (store.ObjectList, error) {

	log.Printf("listing %s", identity)

	var err error
	copt := store.CommonOptionHolder{}
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	res := store.ObjectList{}
	everything := d.PrimaryIndex[identity.Type()]
	if everything == nil {
		return res, nil
	}

	if len(identity.Key()) > 0 {
		return nil, fmt.Errorf("cannot list specific identity")
	}

	for _, v := range everything {
		res = append(res, (*v).Clone())
	}

	// filter results
	// res, err = listFilter(res, copt.Filter)
	// sort results
	res, err = listOrder(res, copt.OrderBy, copt.OrderIncremental)
	if err != nil {
		return nil, err
	}

	// paginate
	return listPagination(res, copt.PageOffset, copt.PageSize), nil
}

func listOrder(list store.ObjectList, ob string, inc bool) (store.ObjectList, error) {
	if len(ob) == 0 {
		return list, nil
	}

	sort.Slice(list, func(p, q int) bool {
		if inc {
			return objectPath(list[p], ob) < objectPath(list[q], ob)
		}
		return objectPath(list[p], ob) > objectPath(list[q], ob)
	})

	return list, nil
}

func listPagination(list store.ObjectList, offset int, size int) store.ObjectList {
	lr := len(list)

	if size == 0 {
		size = lr
	}

	tl := offset
	tr := offset + size
	if lr < tr {
		tr = lr
	}

	if tr <= tl {
		return store.ObjectList{}
	}

	return list[tl:tr]
}

func objectPath(obj store.Object, path string) string {
	data, _ := json.Marshal(obj)
	jsn, err := gabs.ParseJSON(data)
	if err != nil {
		log.Panic(err)
		return ""
	}
	ret := jsn.Path(path).String()
	log.Printf("path %s val %s", path, ret)
	return ret
}
