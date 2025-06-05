package domain

import "time"

// ドキュメント構造体：DynamoDBに保存されるドキュメントを表現
type Document struct {
	// ドキュメントの一意識別子
	ID string `json:"id"`
	// ドキュメントの種類（例：課題、テスト、教材など）
	Type string `json:"type" validate:"required"`
	// ドキュメントの実際のデータ（JSONオブジェクト）
	Data map[string]interface{} `json:"data" validate:"required"`
	// ドキュメントの作成日時
	CreatedAt time.Time `json:"created_at"`
	// ドキュメントの最終更新日時
	UpdatedAt time.Time `json:"updated_at"`
}

// ドキュメント作成リクエスト構造体：新規ドキュメント作成時に使用
type DocumentCreate struct {
	// ドキュメントの種類（必須）
	Type string `json:"type" validate:"required"`
	// ドキュメントのデータ（必須）
	Data map[string]interface{} `json:"data" validate:"required"`
}

// ドキュメント更新リクエスト構造体：既存ドキュメント更新時に使用
type DocumentUpdate struct {
	// 更新するドキュメントのデータ（必須）
	Data map[string]interface{} `json:"data" validate:"required"`
}
