package mongo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/wazofski/storz/internal/constants"
	"github.com/wazofski/storz/internal/logger"
	"github.com/wazofski/storz/internal/utils"
	"github.com/wazofski/storz/store"
	"github.com/wazofski/storz/store/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

var log = logger.Factory("mongo")

type mongoStore struct {
	Schema store.SchemaHolder
	Client *mongo.Client
	Path   string
}

func (d *mongoStore) TestConnection() error {
	const timeout = 10 * time.Second
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	if d.Client != nil {
		if d.Client.Ping(ctx, nil) == nil {
			return nil
		}
	}

	var err error
	d.Client, err = mongo.NewClient(mopt.Client().ApplyURI(d.Path))
	if err != nil {
		return err
	}

	err = d.Client.Connect(ctx)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	err = d.Client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	return nil
}

func Factory(path string) store.Factory {
	return func(schema store.SchemaHolder) (store.Store, error) {
		client := &mongoStore{
			Schema: schema,
			Path:   path,
			Client: nil,
		}

		return client, nil
	}
}

func (d *mongoStore) Create(
	ctx context.Context,
	obj store.Object,
	opt ...options.CreateOption) (store.Object, error) {

	if obj == nil {
		return nil, constants.ErrObjectNil
	}

	log.Printf("create %s", obj.PrimaryKey())

	var err error
	copt := options.CommonOptionHolderFactory()
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

func (d *mongoStore) Update(
	ctx context.Context,
	identity store.ObjectIdentity,
	obj store.Object,
	opt ...options.UpdateOption) (store.Object, error) {

	log.Printf("update %s", identity.Path())

	var err error
	copt := options.CommonOptionHolderFactory()
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

	log.Object("existing", existing)

	err = d.removeIdentity(existing.Metadata().Identity().Path())
	if err != nil {
		log.Printf("%s", err)
	}

	err = d.setIdentity(obj.Metadata().Identity().Path(),
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

	return obj.Clone(), nil
}

func (d *mongoStore) Delete(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.DeleteOption) error {

	log.Printf("delete %s", identity.Path())

	var err error
	copt := options.CommonOptionHolderFactory()
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

func (d *mongoStore) Get(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.GetOption) (store.Object, error) {

	log.Printf("get %s", identity.Path())

	var err error
	copt := options.CommonOptionHolderFactory()
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

func (d *mongoStore) List(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.ListOption) (store.ObjectList, error) {

	log.Printf("list %s", identity.Type())

	var err error
	copt := options.CommonOptionHolderFactory()
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

	query := `SELECT Object FROM Objects
		WHERE Type = ?`

	// pkey filter
	if copt.KeyFilter != nil {
		query = query + fmt.Sprintf(
			" AND Pkey IN ('%s')",
			strings.Join(*copt.KeyFilter, "', '"))
	}

	// prop filter
	if copt.PropFilter != nil {
		query = query + fmt.Sprintf(
			" AND json_extract(Object, '$.%s') = '%s'",
			copt.PropFilter.Key, copt.PropFilter.Value)
	}

	if len(copt.OrderBy) > 0 {
		query = fmt.Sprintf(`SELECT Object
			FROM Objects
			WHERE Type = ?
			ORDER BY json_extract(Object, '$.%s')`, copt.OrderBy)

		if copt.OrderIncremental {
			query = query + " ASC"
		} else {
			query = query + " DESC"
		}
	}

	if copt.PageSize > 0 {
		query = query + fmt.Sprintf(" LIMIT %d", copt.PageSize)
	}

	if copt.PageOffset > 0 {
		query = query + fmt.Sprintf(" OFFSET %d", copt.PageOffset)
	}

	log.Printf(query)

	rows, err := d.DB.Query(query, identity.Type())
	if err != nil {
		return nil, err
	}

	res := d.parseObjectRows(rows, identity.Type())
	rows.Close()

	return res, nil
}

func (d *mongoStore) getIdentity(path string) (string, string, error) {
	row := d.DB.QueryRow("SELECT Pkey, Type FROM IdIndex WHERE Path=?", path)

	var pkey string = ""
	var typ string = ""

	err := row.Scan(&pkey, &typ)
	return pkey, typ, err
}

func (d *mongoStore) setIdentity(path string, pkey string, typ string) error {
	// log.Printf("setting identity %s %s %s", path, pkey, typ)

	query := ""
	_, _, err := d.getIdentity(path)

	if err == nil {
		query = `update IdIndex set Pkey=?, Type=? where Path = ?`
	} else {
		query = `insert into IdIndex (Pkey, Type, Path) values (?, ?, ?)`
	}

	_, err = d.DB.Exec(query, pkey, strings.ToLower(typ), path)

	return err
}

func (d *mongoStore) removeIdentity(path string) error {
	query := "DELETE FROM IdIndex WHERE Path = ?"

	_, err := d.DB.Exec(query, path)
	return err
}

func (d *mongoStore) getObject(pkey string, typ string) (store.Object, error) {
	// log.Printf("getting %s %s", pkey, typ)

	return d.parseObjectRow(
		d.DB.QueryRow("SELECT Object FROM Objects WHERE Pkey=? AND Type=?",
			pkey, strings.ToLower(typ)), typ)
}

func (d *mongoStore) setObject(pkey string, typ string, obj store.Object) error {
	query := ""
	_, err := d.getObject(pkey, typ)
	if err == nil {
		query = `update Objects set Object=@obj where Pkey = @pkey AND Type = @typ`
	} else {
		query = `insert into Objects (Object, Pkey, Type) values (?, ?, ?)`
	}

	data, err := utils.Serialize(obj)

	if err != nil {
		return err
	}

	_, err = d.DB.Exec(query, string(data), pkey, strings.ToLower(typ))
	return err
}

func (d *mongoStore) removeObject(pkey string, typ string) error {
	query := "DELETE FROM Objects WHERE Pkey = ? AND Type = ?"

	_, err := d.DB.Exec(query, pkey, strings.ToLower(typ))

	return err
}

func (d *mongoStore) parseObjectRow(row *sql.Row, typ string) (store.Object, error) {
	var data string = ""

	err := row.Scan(&data)

	if err != nil {
		// log.Fatal(err)
		return nil, err
	}

	return utils.UnmarshalObject([]byte(data), d.Schema, typ)
}

func (d *mongoStore) parseObjectRows(rows *sql.Rows, typ string) store.ObjectList {
	res := store.ObjectList{}
	for rows.Next() {
		var data string = ""
		err := rows.Scan(&data)

		if err != nil {
			log.Fatal(err)
			return nil
		}

		ret, err := utils.UnmarshalObject([]byte(data), d.Schema, typ)
		if err != nil {
			log.Fatal(err)
			return nil
		}

		res = append(res, ret)
	}

	return res
}

func Mongo() {
	/* Connect to my cluster */
	client, err := mongo.NewClient(mopt.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)
	defer cancel()

	/* List databases */
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(databases)

	/* Define my document struct */
	type Post struct {
		Title string `bson:"title,omitempty"`
		Body  string `bson:"body,omitempty"`
	}

	/* Get my collection instance */
	collection := client.Database("blog").Collection("posts")

	/* Insert documents */
	docs := []interface{}{
		bson.D{{Key: "title", Value: "World"}, {Key: "body", Value: "Hello World"}},
		bson.D{{Key: "title", Value: "Mars"}, {Key: "body", Value: "Hello Mars"}},
		bson.D{{Key: "title", Value: "Pluto"}, {Key: "body", Value: "Hello Pluto"}},
	}

	res, insertErr := collection.InsertMany(ctx, docs)
	if insertErr != nil {
		log.Fatal(insertErr)
	}

	fmt.Println(res)

	/* Iterate a cursor and print it */
	cur, currErr := collection.Find(ctx, bson.D{})

	if currErr != nil {
		panic(currErr)
	}

	defer cur.Close(ctx)

	var posts []Post
	if err = cur.All(ctx, &posts); err != nil {
		panic(err)
	}

	fmt.Println(posts)
}
