package store

import "github.com/google/uuid"

type Meta interface {
	Kind() string
	SetKind(string)
	Identity() ObjectIdentity
	SetIdentity(ObjectIdentity)
}

type MetaHolder interface {
	Metadata() Meta
}

type metaWrapper struct {
	Kind_     *string         `json:"kind"`
	Identity_ *ObjectIdentity `json:"identity"`
}

func (m *metaWrapper) Kind() string {
	return *m.Kind_
}

func (m *metaWrapper) Identity() ObjectIdentity {
	return *m.Identity_
}

func (m *metaWrapper) SetKind(kind string) {
	m.Kind_ = &kind
}

func (m *metaWrapper) SetIdentity(identity ObjectIdentity) {
	m.Identity_ = &identity
}

func MetaFactory(kind string) Meta {
	emptyIdentity := ObjectIdentity(uuid.New().String())
	mw := metaWrapper{
		Kind_:     &kind,
		Identity_: &emptyIdentity,
	}

	return &mw
}
