package store

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/aws-samples/serverless-go-demo/types"
)

type DynamoDBStore struct {
	client    *dynamodb.Client
	tableName string
}

var _ types.Store = (*DynamoDBStore)(nil)

func NewDynamoDBStore(ctx context.Context, tableName string) *DynamoDBStore {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return &DynamoDBStore{
		client:    client,
		tableName: tableName,
	}
}

func (d *DynamoDBStore) All(ctx context.Context, next *string) (types.ProductRange, error) {
	productRange := types.ProductRange{
		Products: []types.Product{},
	}

	input := &dynamodb.ScanInput{
		TableName: &d.tableName,
		Limit:     aws.Int32(20),
	}

	if next != nil {
		input.ExclusiveStartKey = map[string]ddbtypes.AttributeValue{
			"id": &ddbtypes.AttributeValueMemberS{Value: *next},
		}
	}

	result, err := d.client.Scan(ctx, input)

	if err != nil {
		return productRange, fmt.Errorf("failed to get items from DynamoDB: %w", err)
	}

	err = attributevalue.UnmarshalListOfMaps(result.Items, &productRange.Products)
	if err != nil {
		return productRange, fmt.Errorf("failed to unmarshal data from DynamoDB: %w", err)
	}

	if len(result.LastEvaluatedKey) > 0 {
		if key, ok := result.LastEvaluatedKey["id"]; ok {
			nextKey := key.(*ddbtypes.AttributeValueMemberS).Value
			productRange.Next = &nextKey
		}
	}

	return productRange, nil
}

func (d *DynamoDBStore) Get(ctx context.Context, id string) (*types.Product, error) {
	response, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &d.tableName,
		Key: map[string]ddbtypes.AttributeValue{
			"id": &ddbtypes.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get item from DynamoDB: %w", err)
	}

	if len(response.Item) == 0 {
		return nil, nil
	}

	product := types.Product{}
	err = attributevalue.UnmarshalMap(response.Item, &product)

	if err != nil {
		return nil, fmt.Errorf("error getting item %w", err)
	}

	return &product, nil
}

func (d *DynamoDBStore) Put(ctx context.Context, product types.Product) error {
	item, err := attributevalue.MarshalMap(&product)
	if err != nil {
		return fmt.Errorf("unable to marshal product: %w", err)
	}

	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &d.tableName,
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("cannot put item: %w", err)
	}

	return nil
}

func (d *DynamoDBStore) Delete(ctx context.Context, id string) error {
	_, err := d.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &d.tableName,
		Key: map[string]ddbtypes.AttributeValue{
			"id": &ddbtypes.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return fmt.Errorf("can't delete item: %w", err)
	}

	return nil
}
