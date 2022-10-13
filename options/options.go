package options

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/wazofski/store"
)

func PropFilter(prop string, val string) store.ListOption {
	return listOption{
		Function: func(options store.OptionHolder) error {
			commonOptions := options.CommonOptions()
			if commonOptions.Filter != nil {
				return errors.New("filter option has already been set")
			}

			commonOptions.Filter = &store.Filter{
				Key:   prop,
				Value: val,
			}
			opstr, _ := json.Marshal(*commonOptions.Filter)
			log.Printf("filter option %s", string(opstr))

			return nil
		},
	}
}

func PageSize(ps int) store.ListOption {
	return listOption{
		Function: func(options store.OptionHolder) error {
			commonOptions := options.CommonOptions()
			if commonOptions.PageSize > 0 {
				return errors.New("page size option has already been set")
			}
			commonOptions.PageSize = ps
			log.Printf("pagination size option %d", ps)
			return nil
		},
	}
}

func PageOffset(po int) store.ListOption {
	return listOption{
		Function: func(options store.OptionHolder) error {
			commonOptions := options.CommonOptions()
			if commonOptions.PageOffset > 0 {
				return errors.New("page offset option has already been set")
			}
			commonOptions.PageOffset = po
			log.Printf("pagination offset option %d", po)
			return nil
		},
	}
}

func OrderBy(field string) store.ListOption {
	return listOption{
		Function: func(options store.OptionHolder) error {
			commonOptions := options.CommonOptions()
			if len(commonOptions.OrderBy) > 0 {
				return errors.New("order by option has already been set")
			}
			commonOptions.OrderBy = field
			log.Printf("order by option: %s", field)
			return nil
		},
	}
}

func OrderIncremental(val bool) store.ListOption {
	return listOption{
		Function: func(options store.OptionHolder) error {
			commonOptions := options.CommonOptions()
			if commonOptions.OrderIncremental {
				return errors.New("order incremental option has already been set")
			}
			commonOptions.OrderIncremental = val
			log.Printf("order incremental option %v", val)
			return nil
		},
	}
}

type listOption struct {
	Function store.OptionFunction
}

func (d listOption) GetListOption() store.Option {
	return d
}

func (d listOption) ApplyFunction() store.OptionFunction {
	return d.Function
}
