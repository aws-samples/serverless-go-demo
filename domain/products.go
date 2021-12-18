package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aws-samples/serverless-go-demo/types"
)

var (
	ErrJsonUnmarshal     = errors.New("failed to parse product from request body")
	ErrProductIdMismatch = errors.New("product ID in path does not match product ID in body")
)

type Products struct {
	store types.Store
}

func NewProductsDomain(s types.Store) *Products {
	return &Products{
		store: s,
	}
}

func (d *Products) GetProduct(ctx context.Context, id string) (*types.Product, error) {
	product, err := d.store.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return product, nil
}

func (d *Products) AllProducts(ctx context.Context, next *string) (types.ProductRange, error) {
	if next != nil && strings.TrimSpace(*next) == "" {
		next = nil
	}

	productRange, err := d.store.All(ctx, next)
	if err != nil {
		return productRange, fmt.Errorf("%w", err)
	}

	return productRange, nil
}

func (d *Products) PutProduct(ctx context.Context, id string, body []byte) (*types.Product, error) {
	product := types.Product{}
	if err := json.Unmarshal(body, &product); err != nil {
		return nil, fmt.Errorf("%w", ErrJsonUnmarshal)
	}

	if product.Id != id {
		return nil, fmt.Errorf("%w", ErrProductIdMismatch)
	}

	err := d.store.Put(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &product, nil
}

func (d *Products) DeleteProduct(ctx context.Context, id string) error {
	err := d.store.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
