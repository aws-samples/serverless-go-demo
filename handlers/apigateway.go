package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/aws-samples/serverless-go-demo/domain"

	"github.com/aws/aws-lambda-go/events"
)

type APIGatewayV2Handler struct {
	products *domain.Products
}

func NewAPIGatewayV2Handler(d *domain.Products) *APIGatewayV2Handler {
	return &APIGatewayV2Handler{
		products: d,
	}
}

func (l *APIGatewayV2Handler) AllHandler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	next := event.QueryStringParameters["next"]

	productRange, err := l.products.AllProducts(ctx, &next)
	if err != nil {
		return errResponse(http.StatusInternalServerError, err.Error()), nil
	}

	return response(http.StatusOK, productRange), nil
}

func (l *APIGatewayV2Handler) GetHandler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	id, ok := event.PathParameters["id"]
	if !ok {
		return errResponse(http.StatusBadRequest, "missing 'id' parameter in path"), nil
	}

	product, err := l.products.GetProduct(ctx, id)

	if err != nil {
		return errResponse(http.StatusInternalServerError, err.Error()), nil
	}
	if product == nil {
		return errResponse(http.StatusNotFound, "product not found"), nil
	} else {
		return response(http.StatusOK, product), nil
	}
}

func (l *APIGatewayV2Handler) PutHandler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	id, ok := event.PathParameters["id"]
	if !ok {
		return errResponse(http.StatusBadRequest, "missing 'id' parameter in path"), nil
	}

	if strings.TrimSpace(event.Body) == "" {
		return errResponse(http.StatusBadRequest, "empty request body"), nil
	}

	product, err := l.products.PutProduct(ctx, id, []byte(event.Body))
	if err != nil {
		if errors.Is(err, domain.ErrJsonUnmarshal) || errors.Is(err, domain.ErrProductIdMismatch) {
			return errResponse(http.StatusBadRequest, err.Error()), nil
		} else {
			return errResponse(http.StatusInternalServerError, err.Error()), nil
		}
	}

	return response(http.StatusCreated, product), nil
}

func (l *APIGatewayV2Handler) DeleteHandler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	id, ok := event.PathParameters["id"]
	if !ok {
		return errResponse(http.StatusBadRequest, "missing 'id' parameter in path"), nil
	}

	err := l.products.DeleteProduct(ctx, id)
	if err != nil {
		return errResponse(http.StatusInternalServerError, err.Error()), nil
	}

	return response(http.StatusOK, nil), nil
}

func response(code int, object interface{}) events.APIGatewayV2HTTPResponse {
	marshalled, err := json.Marshal(object)
	if err != nil {
		return errResponse(http.StatusInternalServerError, err.Error())
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: code,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            string(marshalled),
		IsBase64Encoded: false,
	}
}

func errResponse(status int, body string) events.APIGatewayV2HTTPResponse {
	message := map[string]string{
		"message": body,
	}

	messageBytes, _ := json.Marshal(&message)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(messageBytes),
	}
}
