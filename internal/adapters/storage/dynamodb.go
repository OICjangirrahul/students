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

// DynamoDBストレージ構造体：DynamoDBを使用したドキュメント操作を実装
type DynamoDBStorage struct {
	// DynamoDBクライアント
	client *dynamodb.Client
	// DynamoDBテーブル名
	tableName string
}

// 新しいDynamoDBストレージインスタンスを作成する関数
func NewDynamoDBStorage(client *dynamodb.Client, tableName string) *DynamoDBStorage {
	return &DynamoDBStorage{
		client:    client,
		tableName: tableName,
	}
}

// 新しいドキュメントを作成する
// ドキュメントデータを受け取り、一意のIDを生成して保存し、作成されたドキュメントを返す
func (d *DynamoDBStorage) Create(ctx context.Context, doc *domain.DocumentCreate) (*domain.Document, error) {
	// 現在時刻を取得
	now := time.Now()
	// ドキュメントを作成
	document := &domain.Document{
		ID:        uuid.New().String(),
		Type:      doc.Type,
		Data:      doc.Data,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// ドキュメントをDynamoDB形式に変換
	item, err := attributevalue.MarshalMap(document)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal document: %w", err)
	}

	// DynamoDBにドキュメントを保存
	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      item,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	return document, nil
}

// 指定されたIDのドキュメントを取得する
func (d *DynamoDBStorage) Get(ctx context.Context, id string) (*domain.Document, error) {
	// DynamoDBからドキュメントを取得
	result, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	// ドキュメントが存在しない場合はエラー
	if result.Item == nil {
		return nil, fmt.Errorf("document not found")
	}

	// DynamoDB形式からドキュメントに変換
	var document domain.Document
	err = attributevalue.UnmarshalMap(result.Item, &document)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal document: %w", err)
	}

	return &document, nil
}

// 指定されたIDのドキュメントを更新する
// ドキュメントIDと更新データを受け取り、ドキュメントを更新して更新後のドキュメントを返す
func (d *DynamoDBStorage) Update(ctx context.Context, id string, update *domain.DocumentUpdate) (*domain.Document, error) {
	// 更新式とパラメータを設定
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

	// DynamoDBのドキュメントを更新
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

	// 更新後のドキュメントをDynamoDB形式から変換
	var document domain.Document
	err = attributevalue.UnmarshalMap(result.Attributes, &document)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal updated document: %w", err)
	}

	return &document, nil
}

// 指定されたIDのドキュメントを削除する
func (d *DynamoDBStorage) Delete(ctx context.Context, id string) error {
	// DynamoDBからドキュメントを削除
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

// 指定されたタイプのドキュメント一覧を取得する
func (d *DynamoDBStorage) List(ctx context.Context, docType string) ([]domain.Document, error) {
	// スキャン条件を設定
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

	// DynamoDBからドキュメント一覧を取得
	result, err := d.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}

	// DynamoDB形式からドキュメントリストに変換
	var documents []domain.Document
	err = attributevalue.UnmarshalListOfMaps(result.Items, &documents)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal documents: %w", err)
	}

	return documents, nil
}
