package registry

import (
	"context"
	"fmt"
)

func NewRegistry(name string, options ...OptionRegistry) (registry Registry, err error) {
	if name == "etcd" {
		registry = new(EtcdRegistry)
		err = registry.Init(context.Background(), options...)
	} else {
		err = fmt.Errorf("there is no registry object consist with the given name")
	}

	return
}
