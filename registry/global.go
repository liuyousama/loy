package registry

import (
	"time"
)

var optionsGlobal = &OptionsGlobal{}

type OptionsGlobal struct {
	RegistryPlugin            string
	RegistryEndpoints         []string
	RegistryTimeout           time.Duration
	RegistryHeartBeatDuration int
	RegistryRootPath          string
	CurrentRegistry           Registry
}

func GetCurrentRegistry() Registry {
	return optionsGlobal.CurrentRegistry
}

func Config(options ...OptionGlobal) error {
	for _, opt := range options {
		opt(optionsGlobal)
	}

	if optionsGlobal.RegistryPlugin == "" {
		optionsGlobal.RegistryPlugin = "etcd"
	}

	rgt, err := NewRegistry(optionsGlobal.RegistryPlugin,
		WithEndpoints(optionsGlobal.RegistryEndpoints),
		WithRootPath(optionsGlobal.RegistryRootPath),
		WithTimeout(optionsGlobal.RegistryTimeout))
	if err != nil {
		return err
	}

	optionsGlobal.CurrentRegistry = rgt

	return nil
}

type OptionGlobal func(config *OptionsGlobal)

func WithRegistry(name string) OptionGlobal {
	return func(options *OptionsGlobal) {
		options.RegistryPlugin = name
	}
}

func WithRegistryEndpoints(endpoints []string) OptionGlobal {
	return func(options *OptionsGlobal) {
		options.RegistryEndpoints = endpoints
	}
}

func WithRegistryRootPath(path string) OptionGlobal {
	return func(options *OptionsGlobal) {
		options.RegistryRootPath = path
	}
}
