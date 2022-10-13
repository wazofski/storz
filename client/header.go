package client

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/wazofski/store"
)

type restHeaderOption struct {
	Function store.OptionFunction
}

type headerOption interface {
	store.Option
	store.GetOption
	store.CreateOption
	store.UpdateOption
	store.DeleteOption
	store.ListOption
}

func Header(key string, val string) headerOption {
	return restHeaderOption{
		Function: func(options store.OptionHolder) error {
			restOpts, ok := options.(*restOptions)
			if !ok {
				return errors.New("cannot apply REST specific header option")
			}
			if len(strings.Split(key, " ")) > 1 {
				return fmt.Errorf("invalid header name [%s]", key)
			}
			restOpts.Headers[key] = val
			log.Printf("header option %s%s: [%s]", strings.ToUpper(key[:1]), key[1:], val)
			return nil
		},
	}
}

func (d restHeaderOption) ApplyFunction() store.OptionFunction {
	return d.Function
}
func (d restHeaderOption) GetCreateOption() store.Option {
	return d
}
func (d restHeaderOption) GetDeleteOption() store.Option {
	return d
}
func (d restHeaderOption) GetGetOption() store.Option {
	return d
}
func (d restHeaderOption) GetPatchOption() store.Option {
	return d
}
func (d restHeaderOption) GetUpdateOption() store.Option {
	return d
}
func (d restHeaderOption) GetListOption() store.Option {
	return d
}
func (d restHeaderOption) GetHeaderOption() store.Option {
	return d
}
