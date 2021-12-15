package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws-samples/serverless-go-demo/bus"
	"github.com/aws-samples/serverless-go-demo/domain"
	"github.com/aws-samples/serverless-go-demo/handlers"
)

func main() {
	eventBusName, ok := os.LookupEnv("EVENT_BUS_NAME")
	if !ok {
		panic("Need EVENT_BUS_NAME environment variable")
	}

	store := bus.NewEventBridgeBus(context.TODO(), eventBusName)
	domain := domain.NewProductsStream(store)
	handler := handlers.NewDynamoDBEventHandler(domain)
	lambda.Start(handler.StreamHandler)
}
