package store

import (
	"context"
	"fmt"
	"log"
	"strings"
)

type Object interface {
	MetaHolder
	Clone() Object
	UnmarshalJSON(data []byte) error
	PrimaryKey() string
}

type ObjectList []Object
type ObjectIdentity string

func (o ObjectIdentity) Path() string {
	if strings.Index(string(o), "/") > 0 {
		return string(o)
	}

	return fmt.Sprintf("id/%s", o)
}

func (o ObjectIdentity) Type() string {
	tokens := strings.Split(o.Path(), "/")
	return tokens[0]
}

func (o ObjectIdentity) Key() string {
	tokens := strings.Split(o.Path(), "/")
	if len(tokens) > 1 {
		return tokens[1]
	}
	return ""
}

type Store interface {
	Get(context.Context, ObjectIdentity, ...GetOption) (Object, error)
	List(context.Context, ObjectIdentity, ...ListOption) (ObjectList, error)
	Create(context.Context, Object, ...CreateOption) (Object, error)
	Delete(context.Context, ObjectIdentity, ...DeleteOption) error
	Update(context.Context, ObjectIdentity, Object, ...UpdateOption) (Object, error)
}

type SchemaHolder interface {
	ObjectForKind(kind string) Object
	ObjectMethods() map[string][]string
}

type Factory func(schema SchemaHolder) (Store, error)

func New(schema SchemaHolder, factory Factory) Store {
	store, err := factory(schema)
	if err != nil {
		log.Fatalln(err)
	}
	return store
}
