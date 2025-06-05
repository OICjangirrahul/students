//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/OICjangirrahul/students/internal/adapters/http"
	"github.com/OICjangirrahul/students/internal/adapters/repositories"
	"github.com/OICjangirrahul/students/internal/core/ports"
	"github.com/OICjangirrahul/students/internal/core/services"
	"github.com/google/wire"
)

// データベース依存関係セット：データベース接続を提供
var dbSet = wire.NewSet(repositories.NewDB)

// 学生リポジトリ依存関係セット：学生データの永続化を担当
var studentRepositorySet = wire.NewSet(
	repositories.NewStudentRepository,
	wire.Bind(new(ports.StudentRepository), new(*repositories.StudentRepository)),
)

// 教師リポジトリ依存関係セット：教師データの永続化を担当
var teacherRepositorySet = wire.NewSet(
	repositories.NewTeacherRepository,
	wire.Bind(new(ports.TeacherRepository), new(*repositories.TeacherRepository)),
)

// 学生サービス依存関係セット：学生に関するビジネスロジックを提供
var studentServiceSet = wire.NewSet(
	services.NewStudentService,
)

// 教師サービス依存関係セット：教師に関するビジネスロジックを提供
var teacherServiceSet = wire.NewSet(
	services.NewTeacherService,
)

// ハンドラー構造体：HTTPハンドラーをまとめて管理
type Handlers struct {
	// 学生関連のHTTPハンドラー
	Student *http.StudentHandler
	// 教師関連のHTTPハンドラー
	Teacher *http.TeacherHandler
}

// ハンドラーを初期化する
// 依存関係の注入を行い、必要なコンポーネントを組み立てる
func InitializeHandlers() (*Handlers, error) {
	wire.Build(
		dbSet,
		studentRepositorySet,
		teacherRepositorySet,
		studentServiceSet,
		teacherServiceSet,
		http.NewStudentHandler,
		http.NewTeacherHandler,
		wire.Struct(new(Handlers), "*"),
	)
	return nil, nil
}
