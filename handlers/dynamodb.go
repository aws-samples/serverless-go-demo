package handlers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws-samples/serverless-go-demo/domain"
	"github.com/aws-samples/serverless-go-demo/types"
	"github.com/aws/aws-lambda-go/events"
)

type DynamoDBEventHandler struct {
	productStream *domain.ProductsStream
}

// Can be deleted when this is merged: https://github.com/aws/aws-lambda-go/pull/410/files

type StreamsEventResponse struct {
	BatchItemFailures []BatchItemFailure `json:"batchItemFailures"`
}

type BatchItemFailure struct {
	ItemIdentifier string `json:"itemIdentifier"`
}

func NewDynamoDBEventHandler(p *domain.ProductsStream) *DynamoDBEventHandler {
	return &DynamoDBEventHandler{
		productStream: p,
	}
}

func (d *DynamoDBEventHandler) StreamHandler(ctx context.Context, event events.DynamoDBEvent) (StreamsEventResponse, error) {
	internalEvents := make([]types.Event, len(event.Records))
	for i, ddbEvent := range event.Records {
		internalEvents[i] = eventFromDynamoDBRecord(ddbEvent)
	}

	failedEvents, err := d.productStream.Publish(ctx, internalEvents)
	if err != nil {
		log.Fatalf("totally failed to publish: %v", err)
		return StreamsEventResponse{}, err
	}

	if len(failedEvents) > 0 {
		itemFailures := make([]BatchItemFailure, len(failedEvents))
		for i, failedItem := range failedEvents {
			itemFailures[i] = BatchItemFailure{ItemIdentifier: failedItem.Resources[0]}
		}

		return StreamsEventResponse{BatchItemFailures: itemFailures}, nil
	}

	return StreamsEventResponse{}, nil
}

func eventFromDynamoDBRecord(record events.DynamoDBEventRecord) types.Event {
	change, err := json.Marshal(record.Change)
	if err != nil {
		log.Fatalf("cannot unmarshal dynamodb record change: %s", err)
	}

	detailType := ""
	switch record.EventName {
	case string(events.DynamoDBOperationTypeInsert):
		detailType = "ProductCreated"
	case string(events.DynamoDBOperationTypeModify):
		detailType = "ProductUpdated"
	case string(events.DynamoDBOperationTypeRemove):
		detailType = "ProductDelected"
	}

	return types.Event{
		Source:     "serverless-go-demo",
		Detail:     string(change),
		DetailType: detailType,
		Resources:  []string{record.EventID},
	}
}
