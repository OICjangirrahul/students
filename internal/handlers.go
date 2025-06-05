package internal

import (
	"github.com/OICjangirrahul/students/internal/adapters/http"
	"github.com/OICjangirrahul/students/internal/adapters/repositories"
	"github.com/OICjangirrahul/students/internal/adapters/storage"
	"github.com/OICjangirrahul/students/internal/config"
	"github.com/OICjangirrahul/students/internal/core/services"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Handlers contains all HTTP handlers for the application
type AppHandlers struct {
	Student *http.StudentHandler
	Teacher *http.TeacherHandler
	Storage *http.StorageHandler
}

// InitializeAppHandlers initializes all handlers with their required dependencies
func InitializeAppHandlers(cfg *config.Config) (*AppHandlers, error) {
	db, err := repositories.NewDB(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	studentRepo := repositories.NewStudentRepository(db, cfg)
	teacherRepo := repositories.NewTeacherRepository(db, cfg)

	// Initialize services
	studentService := services.NewStudentService(studentRepo)
	teacherService := services.NewTeacherService(teacherRepo)

	// Initialize AWS clients
	awsCfg, err := config.LoadAWSConfig(cfg)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsCfg)
	dynamoClient := dynamodb.NewFromConfig(awsCfg)

	// Initialize storage services
	fileStorage := storage.NewS3Storage(s3Client, cfg.AWS.S3Bucket)
	documentStorage := storage.NewDynamoDBStorage(dynamoClient, cfg.AWS.DynamoTable)

	// Initialize handlers
	return &AppHandlers{
		Student: http.NewStudentHandler(studentService),
		Teacher: http.NewTeacherHandler(teacherService),
		Storage: http.NewStorageHandler(fileStorage, documentStorage),
	}, nil
}
