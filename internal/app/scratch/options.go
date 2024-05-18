package scratch

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
)

type Options struct {
	ServeMuxOptions []runtime.ServeMuxOption
}

// Option is a configuration callback
type Option interface {
	Apply(*Options) error
}

type optionFn func(*Options) error

func (fn optionFn) Apply(opts *Options) error {
	return fn(opts)
}

func evaluateOptions(opts []Option) (*Options, error) {
	o := &Options{}
	for _, op := range opts {
		err := op.Apply(o)
		if err != nil {
			return nil, errors.Wrap(err, "invalid option")
		}
	}
	return o, nil
}

func WithServeMuxOptions(options ...runtime.ServeMuxOption) Option {
	return optionFn(func(opts *Options) error {
		opts.ServeMuxOptions = options
		return nil
	})
}
