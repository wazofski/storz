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
