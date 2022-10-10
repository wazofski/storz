package store

import (
	"context"
)

type Object interface {
	Metadata() Meta
	Serialize() []byte
}

type ObjectWrapper struct {
	Metadata Meta
}

type ObjectList []Object
type ObjectKind string
type ObjectIdentity string

type Meta interface {
	Kind() ObjectKind
	Identity() ObjectIdentity
}

type metaWrapper struct {
	kind     ObjectKind
	identity ObjectIdentity
}

func (m *metaWrapper) Kind() ObjectKind {
	return m.kind
}

func (m *metaWrapper) Identity() ObjectIdentity {
	return m.identity
}

func ObjectWrapperFactory(kind ObjectKind) ObjectWrapper {
	return ObjectWrapper{
		Metadata: &metaWrapper{
			kind:     kind,
			identity: "",
		},
	}
}

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
