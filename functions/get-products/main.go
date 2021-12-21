package main

import (
	"context"
	"os"

	"github.com/aws-samples/serverless-go-demo/domain"
	"github.com/aws-samples/serverless-go-demo/handlers"
	"github.com/aws-samples/serverless-go-demo/store"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-xray-sdk-go/xray"
)

func main() {
	tableName, ok := os.LookupEnv("TABLE")
	if !ok {
		panic("Need TABLE environment variable")
	}

	xray.Configure(xray.Config{})

	dynamodb := store.NewDynamoDBStore(context.TODO(), tableName)
	domain := domain.NewProductsDomain(dynamodb)
	handler := handlers.NewAPIGatewayV2Handler(domain)
	lambda.Start(handler.AllHandler)
}
