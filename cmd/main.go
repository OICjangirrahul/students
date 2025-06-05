// 学生・教師管理APIのエントリーポイントを提供するパッケージ
package main

import (
	"log/slog"
	"os"

	_ "github.com/OICjangirrahul/students/docs" // Swaggerドキュメントをインポート
	"github.com/OICjangirrahul/students/internal"
	"github.com/OICjangirrahul/students/internal/config"
	"github.com/OICjangirrahul/students/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Student-Teacher Management API
// @version         1.0
// @description     A Go-based REST API for managing students and teachers.
// @host            localhost:8082
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer ` prefix, e.g. "Bearer abcde12345".

func main() {
	// 設定をロード
	cfg, err := config.LoadConfig("config/local.yaml")
	if err != nil {
		slog.Error("failed to load configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// ハンドラーを初期化
	handlers, err := internal.InitializeAppHandlers(cfg)
	if err != nil {
		slog.Error("failed to initialize handlers", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Ginルーターを初期化
	r := gin.Default()

	// CORSミドルウェアを追加
	r.Use(middleware.CorsMiddleware())

	// API v1グループ
	v1 := r.Group("/api/v1")

	// 教師関連のルート
	teachers := v1.Group("/teachers")
	{
		// 公開ルート（認証不要）
		teachers.POST("", handlers.Teacher.Create())      // 教師アカウント作成
		teachers.POST("/login", handlers.Teacher.Login()) // 教師ログイン

		// 保護されたルート（教師ロールが必要）
		protected := teachers.Group("/:id")
		protected.Use(middleware.AuthMiddleware(cfg))           // JWT認証
		protected.Use(middleware.RoleMiddleware("teacher"))     // 教師ロール確認
		protected.Use(middleware.ResourceOwnershipMiddleware()) // リソース所有権確認
		{
			protected.GET("", handlers.Teacher.GetByID())   // 教師情報取得
			protected.PUT("", handlers.Teacher.Update())    // 教師情報更新
			protected.DELETE("", handlers.Teacher.Delete()) // 教師アカウント削除

			// 学生管理ルート
			studentManagement := protected.Group("/students")
			{
				studentManagement.GET("", handlers.Teacher.GetStudents())               // 担当学生一覧取得
				studentManagement.POST("/:studentId", handlers.Teacher.AssignStudent()) // 学生を担当に追加
			}
		}
	}

	// 学生関連のルート
	students := v1.Group("/students")
	{
		// 公開ルート（認証不要）
		students.POST("", handlers.Student.Create())      // 学生アカウント作成
		students.POST("/login", handlers.Student.Login()) // 学生ログイン

		// 保護されたルート（学生または教師ロールが必要）
		protected := students.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))                  // JWT認証
		protected.Use(middleware.RoleMiddleware("student", "teacher")) // ロール確認
		{
			// 学生リソースへのアクセス制御
			protectedStudent := protected.Group("/:id")
			protectedStudent.Use(middleware.ResourceOwnershipMiddleware()) // リソース所有権確認
			{
				protectedStudent.GET("", handlers.Student.GetByID()) // 学生情報取得
			}
		}
	}

	// ストレージ関連のルート（全て認証が必要）
	storage := v1.Group("")
	storage.Use(middleware.AuthMiddleware(cfg)) // JWT認証
	{
		// ファイル関連のルート（教師のみ）
		files := storage.Group("/files")
		files.Use(middleware.RoleMiddleware("teacher")) // 教師ロール確認
		{
			files.POST("", handlers.Storage.UploadFile()) // ファイルアップロード
			files.GET("", handlers.Storage.ListFiles())   // ファイル一覧取得

			fileManagement := files.Group("/:id")
			{
				fileManagement.GET("", handlers.Storage.DownloadFile())  // ファイルダウンロード
				fileManagement.DELETE("", handlers.Storage.DeleteFile()) // ファイル削除
			}
		}

		// ドキュメント関連のルート（教師のみ）
		documents := storage.Group("/documents")
		documents.Use(middleware.RoleMiddleware("teacher")) // 教師ロール確認
		{
			documents.POST("", handlers.Storage.CreateDocument()) // ドキュメント作成
			documents.GET("", handlers.Storage.ListDocuments())   // ドキュメント一覧取得

			documentManagement := documents.Group("/:id")
			{
				documentManagement.GET("", handlers.Storage.GetDocument())       // ドキュメント取得
				documentManagement.PUT("", handlers.Storage.UpdateDocument())    // ドキュメント更新
				documentManagement.DELETE("", handlers.Storage.DeleteDocument()) // ドキュメント削除
			}
		}
	}

	// Swaggerドキュメント
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// サーバーを起動
	if err := r.Run(cfg.HTTPServer.Addr); err != nil {
		slog.Error("failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
