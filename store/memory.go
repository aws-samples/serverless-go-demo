package store

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws-samples/serverless-go-demo/types"
)

type MemoryStore struct {
	storage map[string]types.Product
}

// Just to make sure MemoryStore implements the Store interface
var _ types.Store = (*MemoryStore)(nil)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		storage: make(map[string]types.Product),
	}
}

func (m *MemoryStore) All(ctx context.Context, next *string) (types.ProductRange, error) {
	productRange := types.ProductRange{
		Products: []types.Product{},
		Next:     aws.String("random next string"),
	}

	for _, v := range m.storage {
		productRange.Products = append(productRange.Products, v)
	}

	return productRange, nil
}

func (m *MemoryStore) Get(ctx context.Context, id string) (*types.Product, error) {
	p, ok := m.storage[id]
	if !ok {
		return nil, nil
	}

	return &p, nil
}

func (m *MemoryStore) Put(ctx context.Context, p types.Product) error {
	m.storage[p.Id] = p

	return nil
}

func (m *MemoryStore) Delete(ctx context.Context, id string) error {
	delete(m.storage, id)

	return nil
}
