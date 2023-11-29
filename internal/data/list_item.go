package data

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ListItem struct {
	ID        string    `json:"id"`
	ListID    string    `json:"listId"`
	Item      string    `json:"item"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func marshalListItem(listItem ListItem) (map[string]types.AttributeValue, error) {
	marshaledItem := map[string]types.AttributeValue{
		"ListID":    &types.AttributeValueMemberS{Value: listItem.ListID},
		"ID":        &types.AttributeValueMemberS{Value: listItem.ID},
		"Item":      &types.AttributeValueMemberS{Value: listItem.Item},
		"Done":      &types.AttributeValueMemberBOOL{Value: listItem.Done},
		"CreatedAt": &types.AttributeValueMemberS{Value: listItem.CreatedAt.Format(time.RFC3339)},
		"UpdatedAt": &types.AttributeValueMemberS{Value: listItem.UpdatedAt.Format(time.RFC3339)},
	}

	return marshaledItem, nil
}

func unmarshalListItem(item map[string]types.AttributeValue) (ListItem, error) {
	listID, ok := item["ListID"].(*types.AttributeValueMemberS)
	if !ok {
		return ListItem{}, errors.New("missing ListID attribute")
	}

	itemID, ok := item["ID"].(*types.AttributeValueMemberS)
	if !ok {
		return ListItem{}, errors.New("missing ID attribute")
	}

	itemName, ok := item["Item"].(*types.AttributeValueMemberS)
	if !ok {
		return ListItem{}, errors.New("missing Item attribute")
	}

	done, ok := item["Done"].(*types.AttributeValueMemberBOOL)
	if !ok {
		return ListItem{}, errors.New("missing Done attribute")
	}

	createdAt, ok := item["CreatedAt"].(*types.AttributeValueMemberS)
	if !ok {
		return ListItem{}, errors.New("missing CreatedAt attribute")
	}

	parsedCreatedAt, err := time.Parse(time.RFC3339, *&createdAt.Value)
	if err != nil {
		return ListItem{}, fmt.Errorf("failed to parse CreatedAt: %w", err)
	}

	updatedAt, ok := item["UpdatedAt"].(*types.AttributeValueMemberS)
	if !ok {
		return ListItem{}, errors.New("missing UpdatedAt attribute")
	}

	parsedUpdatedAt, err := time.Parse(time.RFC3339, *&updatedAt.Value)
	if err != nil {
		return ListItem{}, fmt.Errorf("failed to parse UpdatedAt: %w", err)
	}

	return ListItem{
		ID:        itemID.Value,
		ListID:    listID.Value,
		Item:      itemName.Value,
		Done:      done.Value,
		CreatedAt: parsedCreatedAt,
		UpdatedAt: parsedUpdatedAt,
	}, nil
}
