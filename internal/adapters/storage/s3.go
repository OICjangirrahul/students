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

type S3Storage struct {
	client     *s3.Client
	bucketName string
}

func NewS3Storage(client *s3.Client, bucketName string) *S3Storage {
	return &S3Storage{
		client:     client,
		bucketName: bucketName,
	}
}

func (s *S3Storage) Upload(ctx context.Context, file *domain.FileUpload) (*domain.File, error) {
	id := uuid.New().String()
	key := fmt.Sprintf("%s/%s", file.ContentType, id)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(file.Data),
		ContentType: aws.String(file.ContentType),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

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

func (s *S3Storage) Download(ctx context.Context, id string) (*domain.File, []byte, error) {
	file, err := s.Get(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", file.ContentType, id)),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer result.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(result.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file: %w", err)
	}

	return file, buf.Bytes(), nil
}

func (s *S3Storage) Delete(ctx context.Context, id string) error {
	file, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	_, err = s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", file.ContentType, id)),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func (s *S3Storage) List(ctx context.Context) ([]domain.File, error) {
	result, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

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

func (s *S3Storage) Get(ctx context.Context, id string) (*domain.File, error) {
	result, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file metadata: %w", err)
	}

	return &domain.File{
		ID:          id,
		Size:        *result.ContentLength,
		ContentType: aws.ToString(result.ContentType),
		URL:         fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, id),
		BucketName:  s.bucketName,
		UploadedAt:  *result.LastModified,
	}, nil
}
