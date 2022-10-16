package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/wazofski/store"
	"github.com/wazofski/store/constants"
	"github.com/wazofski/store/logger"

	_ "github.com/mattn/go-sqlite3"
)

var log = logger.New("sql")

type sqlStore struct {
	Schema store.SchemaHolder
	Path   string
}

func Factory(path string) store.Factory {
	return func(schema store.SchemaHolder) (store.Store, error) {
		_, err := sql.Open("sqlite3", path)

		if err != nil {
			return nil, err
		}

		client := &sqlStore{
			Schema: schema,
			Path:   path,
		}

		return client, nil
	}
}

/*

// how to create tables
const create string = `
  CREATE TABLE IF NOT EXISTS activities (
  id INTEGER NOT NULL PRIMARY KEY,
  time DATETIME NOT NULL,
  description TEXT
  );`

 if _, err := db.Exec(create); err != nil {

 // inserts
 res, err := c.db.Exec("INSERT INTO activities VALUES(NULL,?,?);", activity.Time, activity.Description)
 if err != nil {
  return 0, err
 }

 var id int64
 if id, err = res.LastInsertId(); err != nil {
  return 0, err
 }

 // selects
 row, err := c.db.Query("SELECT * FROM activities WHERE id=?", id)

 row := c.db.QueryRow("SELECT id, time, description FROM activities WHERE id=?", id)



  // Parse row into Activity struct
 activity := api.Activity{}
 var err error
 if err = row.Scan(&activity.ID, &activity.Time, &activity.Description); err == sql.ErrNoRows {
  log.Printf("Id not found")
  return api.Activity{}, ErrIDNotFound
 }
 return activity, err


// multiple rows
  rows, err := c.db.Query("SELECT * FROM activities WHERE ID > ? ORDER BY id DESC LIMIT 100", offset)
 if err != nil {
  return nil, err
 }
 defer rows.Close()

 data := []api.Activity{}
 for rows.Next() {
  i := api.Activity{}
  err = rows.Scan(&i.ID, &i.Time, &i.Description)
  if err != nil {
   return nil, err
  }
  data = append(data, i)
 }

*/

func (d *sqlStore) Create(
	ctx context.Context,
	obj store.Object,
	opt ...store.CreateOption) (store.Object, error) {

	if obj == nil {
		return nil, constants.ErrObjectNil
	}

	log.Printf("create %s", obj.PrimaryKey())

	var err error
	copt := store.CommonOptionHolder{}
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	lk := strings.ToLower(obj.Metadata().Kind())
	path := fmt.Sprintf("%s/%s", lk, obj.PrimaryKey())
	existing, _ := d.Get(ctx, store.ObjectIdentity(path))

	if existing != nil {
		return nil, constants.ErrObjectExists
	}

	clone := obj.Clone()

	// log.Printf("creating %s", obj.Metadata().Identity())
	// log.Printf("path %s", obj.Metadata().Identity().Path())

	// d.IdentityIndex[obj.Metadata().Identity().Path()] = &clone
	// if d.PrimaryIndex[lk] == nil {
	// 	d.PrimaryIndex[lk] = make(map[string]*store.Object)
	// }

	// d.PrimaryIndex[lk][obj.PrimaryKey()] = &clone

	return clone.Clone(), nil
}

func (d *sqlStore) Update(
	ctx context.Context,
	identity store.ObjectIdentity,
	obj store.Object,
	opt ...store.UpdateOption) (store.Object, error) {

	log.Printf("update %s", identity.Path())

	var err error
	copt := store.CommonOptionHolder{}
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	if obj == nil {
		return nil, constants.ErrObjectNil
	}

	existing, _ := d.Get(ctx, identity)
	if existing == nil {
		return nil, constants.ErrNoSuchObject
	}

	clone := obj.Clone()

	// d.IdentityIndex[obj.Metadata().Identity().Path()] = &clone
	// lk := strings.ToLower(existing.Metadata().Kind())
	// d.PrimaryIndex[lk][existing.PrimaryKey()] = nil

	// lk = strings.ToLower(obj.Metadata().Kind())
	// d.PrimaryIndex[lk][obj.PrimaryKey()] = &clone

	return clone.Clone(), err
}

func (d *sqlStore) Delete(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...store.DeleteOption) error {

	log.Printf("delete %s", identity.Path())

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
		return constants.ErrNoSuchObject
	}

	// d.IdentityIndex[identity.Path()] = nil
	// lk := strings.ToLower(existing.Metadata().Kind())
	// d.PrimaryIndex[lk][existing.PrimaryKey()] = nil

	return nil
}

func (d *sqlStore) Get(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...store.GetOption) (store.Object, error) {

	log.Printf("get %s", identity.Path())

	var err error
	copt := store.CommonOptionHolder{}
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	// log.Printf("...GET identity index size: %d", len(d.IdentityIndex))

	// ret := d.IdentityIndex[identity.Path()]
	// if ret != nil {
	// 	return (*ret).Clone(), nil
	// }

	// tokens := strings.Split(identity.Path(), "/")
	// if len(tokens) == 2 {
	// 	lk := strings.ToLower(tokens[0])
	// 	km := d.PrimaryIndex[lk]
	// 	if km != nil {
	// 		// log.Printf("...GET type index exists with %d records", len(km))
	// 		ret = km[tokens[1]]
	// 		if ret != nil {
	// 			return (*ret).Clone(), nil
	// 		}
	// 	}
	// }

	return nil, constants.ErrNoSuchObject
}

func (d *sqlStore) List(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...store.ListOption) (store.ObjectList, error) {

	log.Printf("list %s", identity.Type())

	var err error
	copt := store.CommonOptionHolder{}
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	res := store.ObjectList{}
	// everything := d.PrimaryIndex[identity.Type()]
	// if everything == nil {
	// 	return res, nil
	// }

	// if len(identity.Key()) > 0 {
	// 	return nil, constants.ErrInvalidPath
	// }

	// for _, v := range everything {
	// 	if v == nil {
	// 		continue
	// 	}
	// 	res = append(res, (*v).Clone())
	// }

	if len(res) > 0 && copt.PropFilter != nil {
		p := objectPath(res[0], copt.PropFilter.Key)
		if p == "" {
			return nil, constants.ErrInvalidFilter
		}
	}

	// key filter results
	res = listPkeyFilter(res, copt.KeyFilter)
	// filter results
	res = listFilter(res, copt.PropFilter)
	// sort results
	res = listOrder(res, copt.OrderBy, copt.OrderIncremental)
	// paginate
	return listPagination(res, copt.PageOffset, copt.PageSize), nil
}

func listPkeyFilter(list store.ObjectList, filter *store.KeyFilter) store.ObjectList {
	if filter == nil {
		return list
	}

	lookup := make(map[string]bool)
	for _, f := range *filter {
		lookup[f] = true
	}

	res := store.ObjectList{}
	for _, o := range list {
		if lookup[o.PrimaryKey()] {
			res = append(res, o)
		}
	}

	return res
}

func listFilter(list store.ObjectList, filter *store.PropFilter) store.ObjectList {
	if filter == nil {
		return list
	}

	res := store.ObjectList{}
	for _, o := range list {
		path := objectPath(o, filter.Key)

		if filter.Value == path {
			res = append(res, o)
		}
	}

	return res
}

func listOrder(list store.ObjectList, ob string, inc bool) store.ObjectList {
	if len(ob) == 0 {
		return list
	}

	sort.Slice(list, func(p, q int) bool {
		if inc {
			return objectPath(list[p], ob) < objectPath(list[q], ob)
		}
		return objectPath(list[p], ob) > objectPath(list[q], ob)
	})

	return list
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
		log.Fatalln(err)
	}
	if !jsn.Exists(strings.Split(path, ".")...) {
		return ""
	}
	ret := strings.ReplaceAll(jsn.Path(path).String(), "\"", "")
	return ret
}
