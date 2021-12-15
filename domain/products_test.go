//go:build unit
// +build unit

package domain

import (
	"context"
	"testing"

	"github.com/aws-samples/serverless-go-demo/store"
	"github.com/aws-samples/serverless-go-demo/types"
)

func TestGetProductNotFound(t *testing.T) {
	memoryStore := store.NewMemoryStore()
	domain := NewProductsDomain(memoryStore)

	product, err := domain.GetProduct(context.Background(), "1")
	if err != nil {
		t.Errorf("GetProduct returned an error: %w", err)
	}

	if product != nil {
		t.Error("GetProduct returned unexpected product")
	}
}

func TestGetExistingProduct(t *testing.T) {
	ctx := context.Background()

	memoryStore := store.NewMemoryStore()
	memoryStore.Put(ctx, types.Product{
		Id:    "iXR",
		Name:  "iPhone XML",
		Price: 0.123,
	})

	domain := NewProductsDomain(memoryStore)

	product, err := domain.GetProduct(context.Background(), "iXR")
	if err != nil {
		t.Errorf("GetProduct returned an error: %w", err)
	}

	if product == nil {
		t.Errorf("GetProduct returned nil object")
		return
	}

	if product.Name != "iPhone XML" {
		t.Errorf("GetProduct returned wrong product name")
	}

	if product.Price != 0.123 {
		t.Errorf("GetProduct returned wrong price")
	}
}
