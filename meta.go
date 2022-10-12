package store

type Meta interface {
	Kind() ObjectKind
	SetKind(ObjectKind)
	Identity() ObjectIdentity
	SetIdentity(ObjectIdentity)
}

type MetaHolder interface {
	Metadata() Meta
}

type metaWrapper struct {
	Kind_     *ObjectKind     `json:"kind"`
	Identity_ *ObjectIdentity `json:"identity"`
}

func (m *metaWrapper) Kind() ObjectKind {
	return *m.Kind_
}

func (m *metaWrapper) Identity() ObjectIdentity {
	return *m.Identity_
}

func (m *metaWrapper) SetKind(kind ObjectKind) {
	m.Kind_ = &kind
}

func (m *metaWrapper) SetIdentity(identity ObjectIdentity) {
	m.Identity_ = &identity
}

func MetaFactory(kind ObjectKind) Meta {
	emptyIdentity := ObjectIdentity("")
	mw := metaWrapper{
		Kind_:     &kind,
		Identity_: &emptyIdentity,
	}

	return &mw
}
