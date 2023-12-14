package data

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoListItemService struct {
	client        *dynamodb.Client
	listItemTable string
	listTable     string
}

func NewDynamoListService(client *dynamodb.Client, listItemTable string, listTable string) *DynamoListItemService {
	return &DynamoListItemService{client: client, listItemTable: listItemTable, listTable: listTable}
}

func (s *DynamoListItemService) IsValidList(ctx context.Context, listId string) error {
	keyEx := expression.Key("ListID").Equal(expression.Value(listId))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		log.Printf("Couldn't build expression for query. Here's why: %v\n", err)
		return err
	}
	input := &dynamodb.QueryInput{
		TableName:                 aws.String(s.listTable),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}
	result, err := s.client.Query(ctx, input)
	if err != nil {
		return err
	}
	if result.Count != 1 {
		return errors.New("list does not exist: " + listId)
	}
	return nil
}

func (s *DynamoListItemService) GetListItems(ctx context.Context, listID string) ([]ListItem, error) {
	filtEx := expression.Name("ListID").Equal(expression.Value(listID))
	expr, err := expression.NewBuilder().WithFilter(filtEx).Build()
	if err != nil {
		log.Printf("Couldn't build expression for query. Here's why: %v\n", err)
		return nil, err
	}
	input := &dynamodb.ScanInput{
		TableName:                 aws.String(s.listItemTable),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	}

	result, err := s.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}

	var listItems []ListItem
	for _, item := range result.Items {
		listItem, err := unmarshalListItem(item)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal item: %w", err)
		}

		listItems = append(listItems, listItem)
	}

	return listItems, nil
}

func (s *DynamoListItemService) CreateListItem(ctx context.Context, listItem *ListItem) error {
	now := time.Now()
	listItem.CreatedAt = now
	listItem.UpdatedAt = now

	marshaledItem, err := marshalListItem(*listItem)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(s.listItemTable),
		Item:      marshaledItem,
	}

	_, err = s.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	return nil
}

func (s *DynamoListItemService) UpdateListItem(ctx context.Context, listItem *ListItem) error {
	listItem.UpdatedAt = time.Now()

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(s.listItemTable),
		Key: map[string]types.AttributeValue{
			"ID":     &types.AttributeValueMemberS{Value: listItem.ID},
			"ListID": &types.AttributeValueMemberS{Value: listItem.ListID},
		},
		UpdateExpression: aws.String("SET Done = :done, UpdatedAt = :updatedAt"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":done":      &types.AttributeValueMemberBOOL{Value: listItem.Done},
			":updatedAt": &types.AttributeValueMemberS{Value: listItem.UpdatedAt.Format(time.RFC3339)},
		},
	}

	_, err := s.client.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}
	return nil
}
