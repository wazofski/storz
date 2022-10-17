package sql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/wazofski/store"
	"github.com/wazofski/store/constants"
	"github.com/wazofski/store/logger"
	"github.com/wazofski/store/utils"

	_ "github.com/mattn/go-sqlite3"
)

var log = logger.New("sql")

type sqlStore struct {
	Schema store.SchemaHolder
	Path   string
	DB     *sql.DB
}

func (d *sqlStore) TestConnection() error {
	if d.DB != nil {
		if d.DB.Ping() == nil {
			return nil
		}
	}

	var err error
	d.DB, err = sql.Open("sqlite3", d.Path)
	if err != nil {
		return err
	}

	err = d.DB.Ping()
	if err != nil {
		return err
	}

	return nil
}

func SqliteFactory(path string) store.Factory {
	return func(schema store.SchemaHolder) (store.Store, error) {
		client := &sqlStore{
			Schema: schema,
			Path:   path,
			DB:     nil,
		}

		log.Printf("initialized %s", path)
		return client, nil
	}
}

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

	err = d.TestConnection()
	if err != nil {
		return nil, err
	}

	err = d.setIdentity(
		obj.Metadata().Identity().Path(),
		obj.PrimaryKey(),
		obj.Metadata().Kind())
	if err != nil {
		return nil, err
	}

	err = d.setObject(obj.PrimaryKey(), obj.Metadata().Kind(), obj)
	if err != nil {
		return nil, err
	}

	return obj.Clone(), nil
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

	err = d.TestConnection()
	if err != nil {
		return nil, err
	}

	err = d.setIdentity(existing.Metadata().Identity().Path(),
		obj.PrimaryKey(), obj.Metadata().Kind())

	if err != nil {
		return nil, err
	}

	err = d.removeObject(existing.PrimaryKey(), existing.Metadata().Kind())
	if err != nil {
		return nil, err
	}

	err = d.setObject(obj.PrimaryKey(), obj.Metadata().Kind(), obj)
	if err != nil {
		return nil, err
	}

	return d.Get(ctx, existing.Metadata().Identity())
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

	err = d.TestConnection()
	if err != nil {
		return err
	}

	err = d.removeIdentity(existing.Metadata().Identity().Path())
	if err != nil {
		return err
	}

	return d.removeObject(existing.PrimaryKey(), existing.Metadata().Kind())
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

	err = d.TestConnection()
	if err != nil {
		return nil, err
	}

	err = d.prepareTables()
	if err != nil {
		return nil, err
	}

	pkey, typ, err := d.getIdentity(identity.Path())
	if err == nil {
		return d.getObject(pkey, typ)
	}

	tokens := strings.Split(identity.Path(), "/")
	if len(tokens) == 2 {
		return d.getObject(tokens[1], tokens[0])
	}

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

	err = d.TestConnection()
	if err != nil {
		return nil, err
	}

	err = d.prepareTables()
	if err != nil {
		return nil, err
	}

	res := store.ObjectList{}

	// selects
	// row, err := c.db.Query("SELECT * FROM activities WHERE id=?", id)

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

	// multiple rows
	// rows, err := c.db.Query("SELECT * FROM activities WHERE ID > ? ORDER BY id DESC LIMIT 100", offset)
	// if err != nil {
	//  return nil, err
	// }
	// defer rows.Close()

	// data := []api.Activity{}
	// for rows.Next() {
	//  i := api.Activity{}
	//  err = rows.Scan(&i.ID, &i.Time, &i.Description)
	//  if err != nil {
	//   return nil, err
	//  }
	//  data = append(data, i)
	// }

	// if len(res) > 0 && copt.PropFilter != nil {
	// 	p := objectPath(res[0], copt.PropFilter.Key)
	// 	if p == "" {
	// 		return nil, constants.ErrInvalidFilter
	// 	}
	// }

	// key filter results
	// res = listPkeyFilter(res, copt.KeyFilter)
	// // filter results
	// res = listFilter(res, copt.PropFilter)
	// // sort results
	// res = listOrder(res, copt.OrderBy, copt.OrderIncremental)
	// // paginate
	// return listPagination(res, copt.PageOffset, copt.PageSize), nil

	return res, nil
}

func (d *sqlStore) prepareTables() error {
	// log.Printf("preparing tables")

	create := `
		CREATE TABLE IF NOT EXISTS IdIndex (
		Path VARCHAR(25) NOT NULL PRIMARY KEY,
		Pkey NVARCHAR(50) NOT NULL,
		Type VARCHAR(25) NOT NULL);`

	_, err := d.DB.Exec(create)
	if err != nil {
		return err
	}

	create = `
		CREATE TABLE IF NOT EXISTS Objects (
		Pkey NVARCHAR(50) NOT NULL,
		Type VARCHAR(25) NOT NULL,
		Object JSON,
		PRIMARY KEY (Pkey,Type));`

	_, err = d.DB.Exec(create)
	if err != nil {
		return err
	}

	return nil
}

func (d *sqlStore) getIdentity(path string) (string, string, error) {
	row := d.DB.QueryRow("SELECT Pkey, Type FROM IdIndex WHERE Path=?", path)

	var pkey string = ""
	var typ string = ""

	err := row.Scan(&pkey, &typ)
	return pkey, typ, err
}

func (d *sqlStore) setIdentity(path string, pkey string, typ string) error {
	query := ""
	_, _, err := d.getIdentity(path)
	if err == nil {
		query = `update IdIndex set Pkey=@pkey, Type=@typ where Path = @path`
	} else {
		query = `insert into IdIndex (Path, Pkey, Type) values (@path, @pkey, @typ)`
	}

	_, err = d.DB.Exec(query,
		sql.Named("path", path),
		sql.Named("pkey", pkey),
		sql.Named("typ", strings.ToLower(typ)))
	return err
}

func (d *sqlStore) removeIdentity(path string) error {
	query := "DELETE FROM IdIndex WHERE Path = @path"

	_, err := d.DB.Exec(query, sql.Named("path", path))
	return err
}

func (d *sqlStore) getObject(pkey string, typ string) (store.Object, error) {
	// log.Printf("getting %s %s", pkey, typ)

	return d.parseObjectRow(
		d.DB.QueryRow("SELECT Object, Type FROM Objects WHERE Pkey=? AND Type=?",
			pkey, strings.ToLower(typ)))
}

func (d *sqlStore) setObject(pkey string, typ string, obj store.Object) error {
	query := ""
	_, err := d.getObject(pkey, typ)
	if err == nil {
		query = `update Objects set Object=@obj where Pkey = @pkey AND Type = @typ`
	} else {
		query = `insert into Objects (Pkey, Type, Object) values (@pkey, @typ, @obj)`
	}

	data, err := utils.Serialize(obj)

	if err != nil {
		return err
	}

	_, err = d.DB.Exec(query,
		sql.Named("pkey", pkey),
		sql.Named("typ", strings.ToLower(typ)),
		sql.Named("obj", string(data)))

	return err
}

func (d *sqlStore) removeObject(pkey string, typ string) error {
	query := "DELETE FROM Objects WHERE Pkey = @pkey AND Type = @typ"

	_, err := d.DB.Exec(query,
		sql.Named("pkey", pkey),
		sql.Named("typ", strings.ToLower(typ)))

	return err
}

func (d *sqlStore) parseObjectRow(row *sql.Row) (store.Object, error) {
	var typ string = ""
	var data string = ""

	err := row.Scan(&data, &typ)

	if err != nil {
		// log.Fatalln(err)
		return nil, err
	}

	return utils.UnmarshalObject([]byte(data), d.Schema, typ)
}
