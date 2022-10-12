package store

import (
	"context"
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

type Option interface {
	ApplyFunction() OptionFunction
}

type CreateOption interface {
	Option
	GetCreateOption() Option
}

type DeleteOption interface {
	Option
	GetDeleteOption() Option
}

type GetOption interface {
	Option
	GetGetOption() Option
}

type UpdateOption interface {
	Option
	GetUpdateOption() Option
}

type ListOption interface {
	Option
	GetListOption() Option
}

type OptionFunction func(OptionHolder) (OptionHolder, error)

type CommonOptionHolder struct {
	// Filter           core.MatcherOp
	OrderBy          string
	OrderIncremental bool
	PageSize         int
	PageOffset       int
}

type OptionHolder interface {
	CommonOptions() *CommonOptionHolder
}

type Store interface {
	Get(context.Context, ObjectIdentity, ...GetOption) (Object, error)
	List(context.Context, ObjectIdentity, ...ListOption) (ObjectList, error)
	Create(context.Context, Object, ...CreateOption) (Object, error)
	Delete(context.Context, ObjectIdentity, ...DeleteOption) error
	Update(context.Context, ObjectIdentity, Object, ...UpdateOption) (Object, error)
}

// type Factory func(schema core.Schema) (Store, error)

// func New(schema core.Schema, factory Factory) Store {
// 	store, err := factory(schema)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	return store
// }
