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

// var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	// flag.Parse()
	// if *cpuprofile != "" {
	// 	f, err := os.Create(*cpuprofile)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	pprof.StartCPUProfile(f)
	// 	defer pprof.StopCPUProfile()
	// }

	xray.Configure(xray.Config{})

	tableName, ok := os.LookupEnv("TABLE")
	if !ok {
		panic("Need TABLE environment variable")
	}

	dynamodb := store.NewDynamoDBStore(context.TODO(), tableName)
	domain := domain.NewProductsDomain(dynamodb)
	handler := handlers.NewAPIGatewayV2Handler(domain)
	lambda.Start(handler.GetHandler)
}
