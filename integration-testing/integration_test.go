//go:build integration
// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws-samples/serverless-go-demo/types"
)

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func getRandomProduct() types.Product {
	return types.Product{
		Id:    randomString(3),
		Name:  randomString(10),
		Price: rand.Float64(),
	}
}

var apiUrl string

func init() {
	_apiUrl, ok := os.LookupEnv("API_URL")
	if !ok {
		panic("Can't find API_URL environment variable")
	}

	apiUrl = _apiUrl
}

func TestFlow(t *testing.T) {
	client := &http.Client{}
	product := getRandomProduct()

	// Put new product
	log.Println("PUT new product")
	payload, err := json.Marshal(product)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", apiUrl, product.Id), bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create product. Got response code %d", resp.StatusCode)
	}

	// Get product
	log.Println("GET product")
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", apiUrl, product.Id), nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to get product. Got response code %d", resp.StatusCode)
	}

	apiProduct := types.Product{}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &apiProduct)

	if apiProduct.Id != product.Id {
		t.Fatalf("API product ID is different from our own")
	}

	if apiProduct.Name != product.Name {
		t.Fatalf("API product Name is different from our own")
	}

	if apiProduct.Price != product.Price {
		t.Fatalf("API product Price is different from our own")
	}

	// Get all products
	log.Println("GET all products")

	req, err = http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to get all products. Got response code %d", resp.StatusCode)
	}

	products := types.ProductRange{}
	body, _ = io.ReadAll(resp.Body)
	json.Unmarshal(body, &products)

	if len(products.Products) < 1 {
		t.Fatalf("Failed to get all products. Only got %d", len(products.Products))
	}

	// Delete product
	log.Println("DELETE product")

	req, err = http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", apiUrl, product.Id), nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to delete product. Got response code %d", resp.StatusCode)
	}

	// Get product again
	log.Println("GET product again")
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", apiUrl, product.Id), nil)
	if err != nil {
		panic(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatal("Got unexpected product")
	}
}

func TestPutProductWithInvalidId(t *testing.T) {
	client := &http.Client{}

	product := getRandomProduct()
	product.Id = "invalid id"

	log.Println("PUT new product")
	payload, err := json.Marshal(product)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", apiUrl, "not-the-same-id"), bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Should have not created product. Got response code %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "product ID in path does not match product ID in body") {
		t.Fatalf("Wrong body content: %s", string(body))
	}
}

func TestProductEmpty(t *testing.T) {
	client := &http.Client{}

	log.Println("PUT new product")
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", apiUrl, "empty-id"), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Should have not created product. Got response code %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "empty request body") {
		t.Fatalf("Wrong body content: %s", string(body))
	}
}

func TestPutProductInvalidBody(t *testing.T) {
	client := &http.Client{}

	product := getRandomProduct()

	log.Println("PUT new product")
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", apiUrl, product.Id), bytes.NewReader([]byte("invalid body")))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Should have not created product. Got response code %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "failed to parse product from request body") {
		t.Fatalf("Wrong body content: %s", string(body))
	}
}
