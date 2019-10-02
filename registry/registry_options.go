package registry

import "time"

type OptionsRegistry struct {
	Endpoints []string
	Timeout time.Duration
	HeartBeatDuration int
	RootPath string
}

type OptionRegistry func(options *OptionsRegistry)

func WithTimeout(timeout time.Duration) OptionRegistry {
	return func(options *OptionsRegistry) {
		options.Timeout = timeout
	}
}

func WithEndpoints(endpoints []string) OptionRegistry {
	return func(options *OptionsRegistry) {
		options.Endpoints = endpoints
	}
}

func WithRootPath(path string) OptionRegistry {
	return func(options *OptionsRegistry) {
		options.RootPath = path
	}
}