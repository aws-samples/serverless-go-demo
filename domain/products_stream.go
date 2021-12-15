package domain

import (
	"context"
	"fmt"

	"github.com/aws-samples/serverless-go-demo/types"
)

type ProductsStream struct {
	bus types.Bus
}

func NewProductsStream(b types.Bus) *ProductsStream {
	return &ProductsStream{
		bus: b,
	}
}

func (p *ProductsStream) Publish(ctx context.Context, events []types.Event) ([]types.FailedEvent, error) {
	failedEvents, err := p.bus.Put(ctx, events)
	if err != nil {
		return failedEvents, fmt.Errorf("%w", err)
	}

	return failedEvents, nil
}
