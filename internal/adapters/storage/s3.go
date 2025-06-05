package storage

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/OICjangirrahul/students/internal/core/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// S3ストレージ構造体：S3を使用したファイル操作を実装
type S3Storage struct {
	// S3クライアント
	client *s3.Client
	// S3バケット名
	bucketName string
}

// 新しいS3ストレージインスタンスを作成する関数
func NewS3Storage(client *s3.Client, bucketName string) *S3Storage {
	return &S3Storage{
		client:     client,
		bucketName: bucketName,
	}
}

// ファイルをS3にアップロードする
// 一意のIDを生成し、ファイルを保存して、保存されたファイルの情報を返す
func (s *S3Storage) Upload(ctx context.Context, file *domain.FileUpload) (*domain.File, error) {
	// 一意のIDを生成
	id := uuid.New().String()
	key := fmt.Sprintf("%s/%s", file.ContentType, id)

	// S3にファイルをアップロード
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(file.Data),
		ContentType: aws.String(file.ContentType),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// アップロードされたファイルの情報を返す
	return &domain.File{
		ID:          id,
		Name:        file.Name,
		Size:        int64(len(file.Data)),
		ContentType: file.ContentType,
		URL:         fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, key),
		BucketName:  s.bucketName,
		UploadedAt:  time.Now(),
	}, nil
}

// S3からファイルをダウンロードする
// ファイルIDを受け取り、ファイル情報とデータを返す
func (s *S3Storage) Download(ctx context.Context, id string) (*domain.File, []byte, error) {
	// ファイル情報を取得
	file, err := s.Get(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	// S3からファイルをダウンロード
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", file.ContentType, id)),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer result.Body.Close()

	// ファイルデータを読み込む
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(result.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}

	return file, buf.Bytes(), nil
}

// S3からファイルを削除する
// ファイルIDを受け取り、対応するファイルを削除する
func (s *S3Storage) Delete(ctx context.Context, id string) error {
	// ファイル情報を取得
	file, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	// S3からファイルを削除
	_, err = s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", file.ContentType, id)),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// S3に保存されているファイルの一覧を取得する
func (s *S3Storage) List(ctx context.Context) ([]domain.File, error) {
	// S3からファイル一覧を取得
	result, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	// ファイル情報を変換
	files := make([]domain.File, 0, len(result.Contents))
	for _, obj := range result.Contents {
		file := domain.File{
			ID:         aws.ToString(obj.Key),
			Size:       *obj.Size,
			URL:        fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, aws.ToString(obj.Key)),
			BucketName: s.bucketName,
			UploadedAt: *obj.LastModified,
		}
		files = append(files, file)
	}

	return files, nil
}

// 指定されたIDのファイル情報を取得する
func (s *S3Storage) Get(ctx context.Context, id string) (*domain.File, error) {
	// S3からファイルのメタデータを取得
	result, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file metadata: %w", err)
	}

	// ファイル情報を返す
	return &domain.File{
		ID:          id,
		Size:        *result.ContentLength,
		ContentType: aws.ToString(result.ContentType),
		URL:         fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, id),
		BucketName:  s.bucketName,
		UploadedAt:  *result.LastModified,
	}, nil
}
