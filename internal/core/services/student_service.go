package services

import (
	"context"

	"github.com/OICjangirrahul/students/internal/core/domain"
	"github.com/OICjangirrahul/students/internal/core/ports"
)

// 学生サービス構造体：学生に関する業務ロジックを実装
type StudentService struct {
	// 学生リポジトリインターフェース
	repo ports.StudentRepository
}

// 新しい学生サービスインスタンスを作成する
func NewStudentService(repo ports.StudentRepository) *StudentService {
	return &StudentService{
		repo: repo,
	}
}

// 新しい学生を作成する
// 学生データを受け取り、データベースに保存して、作成された学生の情報を返す
func (s *StudentService) Create(ctx context.Context, student *domain.Student) (*domain.Student, error) {
	// 学生を作成
	id, err := s.repo.CreateStudent(student.Name, student.Email, student.Age, student.Password)
	if err != nil {
		return nil, err
	}

	// タイムスタンプを含む完全な学生情報を取得
	return s.repo.GetStudentByID(id)
}

// 指定されたIDの学生を取得する
func (s *StudentService) GetByID(ctx context.Context, id int64) (*domain.Student, error) {
	return s.repo.GetStudentByID(id)
}

// 学生のログイン認証を行う
// メールアドレスとパスワードを検証し、有効な場合はJWTトークンを返す
func (s *StudentService) Login(ctx context.Context, email, password string) (string, error) {
	return s.repo.LoginStudent(email, password)
}
