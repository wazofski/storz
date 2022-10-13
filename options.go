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

type CommonOptionHolder struct {
	// Filter           core.MatcherOp
	OrderBy          string
	OrderIncremental bool
	PageSize         int
	PageOffset       int
}

func (d *CommonOptionHolder) CommonOptions() *CommonOptionHolder {
	return d
}
