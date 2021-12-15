//go:build unit
// +build unit

package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/aws-samples/serverless-go-demo/store"
	"github.com/aws-samples/serverless-go-demo/types"
	"github.com/aws-samples/serverless-go-demo/types/mocks"
	"github.com/golang/mock/gomock"
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

func TestGetFailedStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	store := mocks.NewMockStore(ctrl)
	store.EXPECT().
		Get(ctx, gomock.Eq("1")).
		Return(nil, errors.New("internal error"))

	domain := NewProductsDomain(store)

	product, err := domain.GetProduct(context.Background(), "1")
	if product != nil {
		t.Error("Got unexpected product")
	}

	if err == nil {
		t.Error("Expecting an error to be returned")
		return
	}

	if err.Error() != "internal error" {
		t.Errorf("Got unexpected error: %s", err)
	}
}
