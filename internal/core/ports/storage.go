package ports

import (
	"context"

	"github.com/OICjangirrahul/students/internal/core/domain"
)

// FileStorage defines the interface for file storage operations (S3)
type FileStorage interface {
	Upload(ctx context.Context, file *domain.FileUpload) (*domain.File, error)
	Download(ctx context.Context, id string) (*domain.File, []byte, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.File, error)
	Get(ctx context.Context, id string) (*domain.File, error)
}

// DocumentStorage defines the interface for document storage operations (DynamoDB)
type DocumentStorage interface {
	Create(ctx context.Context, doc *domain.DocumentCreate) (*domain.Document, error)
	Get(ctx context.Context, id string) (*domain.Document, error)
	Update(ctx context.Context, id string, update *domain.DocumentUpdate) (*domain.Document, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, docType string) ([]domain.Document, error)
}
