package mongo

import (
	"context"
	"encoding/json"
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

const collectionName = "objects"
const timeout = 10 * time.Second

type mongoStore struct {
	Schema store.SchemaHolder
	Client *mongo.Client
	Path   string
	DB     string
}

type _Record struct {
	IdPath string      `json:"idpath" bson:"idpath"`
	PkPath string      `json:"pkpath" bson:"pkpath"`
	Type   string      `json:"type" bson:"type"`
	Obj    interface{} `json:"object" bson:"object"`
}

func (d *mongoStore) TestConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

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

	return d.prepare()
}

func (d *mongoStore) prepare() error {
	collection := d.Client.Database(d.DB).Collection(collectionName)
	indexModel := []mongo.IndexModel{
		{
			Keys: bson.M{
				"idpath": 1,
			}, Options: nil,
		},
		{
			Keys: bson.M{
				"pkpath": 1,
			}, Options: nil,
		},
		{
			Keys: bson.M{
				"type": 1,
			}, Options: nil,
		},
	}

	for _, i := range indexModel {
		collection.Indexes().CreateOne(context.Background(), i)
	}

	return nil
}

func Factory(path string, db string) store.Factory {
	return func(schema store.SchemaHolder) (store.Store, error) {
		client := &mongoStore{
			Schema: schema,
			Path:   path,
			DB:     db,
			Client: nil,
		}

		err := client.TestConnection()
		if err != nil {
			return nil, err
		}

		log.Printf("initialized %s %s", path, db)
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

	typ := strings.ToLower(obj.Metadata().Kind())

	collection := d.Client.Database(d.DB).Collection(collectionName)
	_, err = collection.InsertOne(ctx,
		_Record{
			IdPath: obj.Metadata().Identity().Path(),
			PkPath: fmt.Sprintf("%s/%s", typ, obj.PrimaryKey()),
			Type:   typ,
			Obj:    toBSON(obj),
		})

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

	err = d.TestConnection()
	if err != nil {
		return nil, err
	}

	err = d.Delete(ctx, identity)
	if err != nil {
		return nil, err
	}

	return d.Create(ctx, obj)
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

	collection := d.Client.Database(d.DB).Collection(collectionName)
	_, err = collection.DeleteOne(ctx,
		bson.M{
			"idpath": identity.Path(),
			"pkpath": identity.Path(),
		})

	return err
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

	collection := d.Client.Database(d.DB).Collection(collectionName)
	var res bson.M
	collection.FindOne(ctx,
		bson.M{
			"idpath": identity.Path(),
		}).Decode(&res)

	if res != nil {
		return fromBSON(res, d.Schema)
	}

	collection.FindOne(ctx,
		bson.M{
			"pkpath": identity.Path(),
		}).Decode(&res)

	if res != nil {
		return fromBSON(res, d.Schema)
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

	// rows, err := d.DB.Query(query, identity.Type())
	// if err != nil {
	// 	return nil, err
	// }

	// res := d.parseObjectRows(rows, identity.Type())
	// rows.Close()

	return nil, nil
}

func toBSON(obj store.Object) interface{} {
	data, _ := utils.Serialize(obj)
	res := make(map[string]interface{})
	json.Unmarshal(data, &res)
	return res
}

func fromBSON(bs interface{}, schema store.SchemaHolder) (store.Object, error) {
	m := bs.(bson.M)
	if m == nil {
		return nil, fmt.Errorf("invalid bson")
	}

	data, err := json.Marshal(m["object"])
	if err != nil {
		return nil, err
	}

	return utils.UnmarshalObject(data, schema, utils.ObjeectKind(data))

}
