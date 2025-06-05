package ports

import (
	"context"

	"github.com/OICjangirrahul/students/internal/core/domain"
)

// ファイルストレージインターフェース：S3を使用したファイル操作を定義
type FileStorage interface {
	// ファイルをアップロードし、保存されたファイルの情報を返す
	Upload(ctx context.Context, file *domain.FileUpload) (*domain.File, error)
	// 指定されたIDのファイルをダウンロードし、ファイル情報とデータを返す
	Download(ctx context.Context, id string) (*domain.File, []byte, error)
	// 指定されたIDのファイルを削除する
	Delete(ctx context.Context, id string) error
	// 保存されているファイルの一覧を取得する
	List(ctx context.Context) ([]domain.File, error)
	// 指定されたIDのファイル情報を取得する
	Get(ctx context.Context, id string) (*domain.File, error)
}

// ドキュメントストレージインターフェース：DynamoDBを使用したドキュメント操作を定義
type DocumentStorage interface {
	// 新しいドキュメントを作成し、作成されたドキュメントを返す
	Create(ctx context.Context, doc *domain.DocumentCreate) (*domain.Document, error)
	// 指定されたIDのドキュメントを取得する
	Get(ctx context.Context, id string) (*domain.Document, error)
	// 指定されたIDのドキュメントを更新し、更新後のドキュメントを返す
	Update(ctx context.Context, id string, update *domain.DocumentUpdate) (*domain.Document, error)
	// 指定されたIDのドキュメントを削除する
	Delete(ctx context.Context, id string) error
	// 指定されたタイプのドキュメント一覧を取得する
	List(ctx context.Context, docType string) ([]domain.Document, error)
}
