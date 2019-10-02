package registry

import (
	"context"
)

type Registry interface {
	Name() string

	Init(ctx context.Context, options ...OptionRegistry) error

	Register(ctx context.Context, service *Service)

	Withdraw(ctx context.Context, service *Service)

	Discover(ctx context.Context, name string) (*Service, error)
}
