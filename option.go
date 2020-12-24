package registry

import (
	"time"
)

type Options struct {
	Addrs     []string
	Timeout   time.Duration
	Path      string
	HeartBeat int64
}

type Option func(opts *Options)

func WithTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.Timeout = timeout
	}
}

func WithAddrs(addrs []string) Option {
	return func(opts *Options) {
		opts.Addrs = addrs
	}
}

func WithPath(path string) Option {
	return func(opts *Options) {
		opts.Path = path
	}
}

func WithHeartBeat(heartHeat int64) Option {
	return func(opts *Options) {
		opts.HeartBeat = heartHeat
	}
}
