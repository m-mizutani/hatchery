package hatchery

import "context"

type Source interface {
	Read(ctx context.Context) error
}
