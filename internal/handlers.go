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

// アプリケーションハンドラー構造体：全てのHTTPハンドラーを管理
type AppHandlers struct {
	// 学生関連のHTTPハンドラー
	Student *http.StudentHandler
	// 教師関連のHTTPハンドラー
	Teacher *http.TeacherHandler
	// ストレージ関連のHTTPハンドラー
	Storage *http.StorageHandler
}

// アプリケーションハンドラーを初期化する
// 必要な依存関係を注入し、全てのハンドラーを設定する
func InitializeAppHandlers(cfg *config.Config) (*AppHandlers, error) {
	// データベース接続を初期化
	db, err := repositories.NewDB(cfg)
	if err != nil {
		return nil, err
	}

	// リポジトリを初期化
	// データベースとの対話を担当するコンポーネントを作成
	studentRepo := repositories.NewStudentRepository(db, cfg)
	teacherRepo := repositories.NewTeacherRepository(db, cfg)

	// サービスを初期化
	// ビジネスロジックを実装するコンポーネントを作成
	studentService := services.NewStudentService(studentRepo)
	teacherService := services.NewTeacherService(teacherRepo)

	// AWSクライアントを初期化
	// S3とDynamoDBへのアクセスを設定
	awsCfg, err := config.LoadAWSConfig(cfg)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsCfg)
	dynamoClient := dynamodb.NewFromConfig(awsCfg)

	// ストレージサービスを初期化
	// ファイルとドキュメントの保存を担当するコンポーネントを作成
	fileStorage := storage.NewS3Storage(s3Client, cfg.AWS.S3Bucket)
	documentStorage := storage.NewDynamoDBStorage(dynamoClient, cfg.AWS.DynamoTable)

	// ハンドラーを初期化して返す
	// 各種サービスを利用してHTTPリクエストを処理するハンドラーを作成
	return &AppHandlers{
		Student: http.NewStudentHandler(studentService),
		Teacher: http.NewTeacherHandler(teacherService),
		Storage: http.NewStorageHandler(fileStorage, documentStorage),
	}, nil
}
