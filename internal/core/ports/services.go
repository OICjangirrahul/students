package ports

import (
	"context"

	"github.com/OICjangirrahul/students/internal/core/domain"
)

// 学生サービスインターフェース：学生に関する業務ロジックを定義
type StudentService interface {
	// 新しい学生を作成する
	Create(ctx context.Context, student *domain.Student) (*domain.Student, error)
	// 指定されたIDの学生を取得する
	GetByID(ctx context.Context, id int64) (*domain.Student, error)
	// 学生のログイン認証を行い、JWTトークンを返す
	Login(ctx context.Context, email, password string) (string, error)
}

// 教師サービスインターフェース：教師に関する業務ロジックを定義
type TeacherService interface {
	// 新しい教師を作成する
	Create(ctx context.Context, teacher *domain.Teacher) (*domain.Teacher, error)
	// 指定されたIDの教師を取得する
	GetByID(ctx context.Context, id int64) (*domain.Teacher, error)
	// 教師情報を更新する
	Update(ctx context.Context, teacher *domain.Teacher) (*domain.Teacher, error)
	// 指定されたIDの教師を削除する
	Delete(ctx context.Context, id int64) error
	// 教師のログイン認証を行い、JWTトークンを返す
	Login(ctx context.Context, email, password string) (string, error)
	// 教師に学生を割り当てる
	AssignStudent(ctx context.Context, teacherID, studentID int64) error
	// 教師に割り当てられた学生一覧を取得する
	GetStudents(ctx context.Context, teacherID int64) ([]domain.Student, error)
}
