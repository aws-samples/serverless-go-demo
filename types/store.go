package types

import (
	"context"
)

type Store interface {
	All(context.Context, *string) (ProductRange, error)
	Get(context.Context, string) (*Product, error)
	Put(context.Context, Product) error
	Delete(context.Context, string) error
}
