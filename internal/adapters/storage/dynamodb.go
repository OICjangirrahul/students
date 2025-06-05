package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/OICjangirrahul/students/internal/core/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type DynamoDBStorage struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBStorage(client *dynamodb.Client, tableName string) *DynamoDBStorage {
	return &DynamoDBStorage{
		client:    client,
		tableName: tableName,
	}
}

func (d *DynamoDBStorage) Create(ctx context.Context, doc *domain.DocumentCreate) (*domain.Document, error) {
	now := time.Now()
	document := &domain.Document{
		ID:        uuid.New().String(),
		Type:      doc.Type,
		Data:      doc.Data,
		CreatedAt: now,
		UpdatedAt: now,
	}

	item, err := attributevalue.MarshalMap(document)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal document: %w", err)
	}

	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      item,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	return document, nil
}

func (d *DynamoDBStorage) Get(ctx context.Context, id string) (*domain.Document, error) {
	result, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("document not found")
	}

	var document domain.Document
	err = attributevalue.UnmarshalMap(result.Item, &document)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal document: %w", err)
	}

	return &document, nil
}

func (d *DynamoDBStorage) Update(ctx context.Context, id string, update *domain.DocumentUpdate) (*domain.Document, error) {
	updateExpr := "SET #data = :data, #updatedAt = :updatedAt"
	attrNames := map[string]string{
		"#data":      "Data",
		"#updatedAt": "UpdatedAt",
	}
	attrValues, err := attributevalue.MarshalMap(map[string]interface{}{
		":data":      update.Data,
		":updatedAt": time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update values: %w", err)
	}

	result, err := d.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(d.tableName),
		Key:                       map[string]types.AttributeValue{"ID": &types.AttributeValueMemberS{Value: id}},
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeNames:  attrNames,
		ExpressionAttributeValues: attrValues,
		ReturnValues:              types.ReturnValueAllNew,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	var document domain.Document
	err = attributevalue.UnmarshalMap(result.Attributes, &document)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal updated document: %w", err)
	}

	return &document, nil
}

func (d *DynamoDBStorage) Delete(ctx context.Context, id string) error {
	_, err := d.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	return nil
}

func (d *DynamoDBStorage) List(ctx context.Context, docType string) ([]domain.Document, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String(d.tableName),
		FilterExpression: aws.String("#type = :type"),
		ExpressionAttributeNames: map[string]string{
			"#type": "Type",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":type": &types.AttributeValueMemberS{Value: docType},
		},
	}

	result, err := d.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}

	var documents []domain.Document
	err = attributevalue.UnmarshalListOfMaps(result.Items, &documents)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal documents: %w", err)
	}

	return documents, nil
}
