package store

import (
	"context"
	"log"
)

type Object interface {
	MetaHolder
	Clone() Object
}

type SchemaHolder interface {
	ObjectForKind(kind string) Object
	ObjectMethods() map[string][]string
}

type ObjectList []Object
type ObjectIdentity string

type Store interface {
	Get(context.Context, ObjectIdentity, ...GetOption) (Object, error)
	List(context.Context, ObjectIdentity, ...ListOption) (ObjectList, error)
	Create(context.Context, Object, ...CreateOption) (Object, error)
	Delete(context.Context, ObjectIdentity, ...DeleteOption) error
	Update(context.Context, ObjectIdentity, Object, ...UpdateOption) (Object, error)
}

type Factory func(schema SchemaHolder) (Store, error)

func New(schema SchemaHolder, factory Factory) Store {
	store, err := factory(schema)
	if err != nil {
		log.Fatalln(err)
	}
	return store
}
