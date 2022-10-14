package store

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

type OptionHolder interface {
	CommonOptions() *CommonOptionHolder
}

type OptionFunction func(OptionHolder) error

type Filter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CommonOptionHolder struct {
	Filter           *Filter
	OrderBy          string
	OrderIncremental bool
	PageSize         int
	PageOffset       int
}

func (d *CommonOptionHolder) CommonOptions() *CommonOptionHolder {
	return d
}