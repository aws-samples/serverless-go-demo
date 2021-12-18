package bus

import (
	"context"
	"log"
	"math"

	"github.com/aws-samples/serverless-go-demo/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchevents/types"
)

type EventBridgeBus struct {
	client  *cloudwatchevents.Client
	busName string
}

var _ types.Bus = (*EventBridgeBus)(nil)

func NewEventBridgeBus(ctx context.Context, busName string) *EventBridgeBus {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := cloudwatchevents.NewFromConfig(cfg)

	return &EventBridgeBus{
		client:  client,
		busName: busName,
	}
}

func (e *EventBridgeBus) Put(ctx context.Context, events []types.Event) ([]types.FailedEvent, error) {
	failedBatchEvents, err :=
		batchEvents(events, 10, func(batchEvents []types.Event) ([]types.FailedEvent, error) {
			eventBridgeEvents := make([]cloudwatchtypes.PutEventsRequestEntry, len(batchEvents))

			for i, event := range batchEvents {
				eventBridgeEvent := cloudwatchtypes.PutEventsRequestEntry{
					EventBusName: &e.busName,
					Source:       aws.String(event.Source),
					Detail:       aws.String(event.Detail),
					DetailType:   aws.String(event.DetailType),
					Resources:    event.Resources,
				}

				eventBridgeEvents[i] = eventBridgeEvent
			}

			result, err := e.client.PutEvents(ctx, &cloudwatchevents.PutEventsInput{
				Entries: eventBridgeEvents,
			})

			failedEvents := []types.FailedEvent{}
			if err != nil {
				return failedEvents, err
			}

			if result.FailedEntryCount > 0 {
				for i, entry := range result.Entries {
					if entry.EventId != nil {
						continue
					}

					failedEvent := types.FailedEvent{
						Event:          batchEvents[i],
						FailureCode:    *entry.ErrorCode,
						FailureMessage: *entry.ErrorMessage,
					}

					failedEvents = append(failedEvents, failedEvent)
				}
			}

			return failedEvents, nil
		})

	return failedBatchEvents, err
}

func batchEvents(events []types.Event, maxBatchSize uint, batchFn func([]types.Event) ([]types.FailedEvent, error)) ([]types.FailedEvent, error) {
	skip := 0
	recordsAmount := len(events)
	batchAmount := int(math.Ceil(float64(recordsAmount) / float64(maxBatchSize)))

	batchFailedEvents := []types.FailedEvent{}

	for i := 0; i < batchAmount; i++ {
		lowerBound := skip
		upperBound := skip + int(maxBatchSize)

		if upperBound > recordsAmount {
			upperBound = recordsAmount
		}

		batchEvents := events[lowerBound:upperBound]
		skip += int(maxBatchSize)

		failedEvents, err := batchFn(batchEvents)
		if err != nil {
			return batchFailedEvents, err
		}

		if len(failedEvents) > 0 {
			batchFailedEvents = append(batchFailedEvents, failedEvents...)
		}
	}

	return batchFailedEvents, nil
}
