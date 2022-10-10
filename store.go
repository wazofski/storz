package store

type Meta interface {
	Kind() string
	Identity() string
}

type Object interface {
	Meta() Meta
	Serialize() []byte
}

type ObjectList []Object
