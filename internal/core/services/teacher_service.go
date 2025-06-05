package services

import (
	"context"

	"github.com/OICjangirrahul/students/internal/core/domain"
	"github.com/OICjangirrahul/students/internal/core/ports"
)

// 教師サービス構造体：教師に関する業務ロジックを実装
type TeacherService struct {
	// 教師リポジトリインターフェース
	repo ports.TeacherRepository
}

// 新しい教師サービスインスタンスを作成する
func NewTeacherService(repo ports.TeacherRepository) *TeacherService {
	return &TeacherService{
		repo: repo,
	}
}

// 新しい教師を作成する
// 教師データを受け取り、データベースに保存して、作成された教師の情報を返す
func (s *TeacherService) Create(ctx context.Context, teacher *domain.Teacher) (*domain.Teacher, error) {
	// 教師を作成
	id, err := s.repo.CreateTeacher(teacher.Name, teacher.Email, teacher.Password, teacher.Subject)
	if err != nil {
		return nil, err
	}

	// タイムスタンプを含む完全な教師情報を取得
	return s.repo.GetTeacherByID(id)
}

// 指定されたIDの教師を取得する
func (s *TeacherService) GetByID(ctx context.Context, id int64) (*domain.Teacher, error) {
	return s.repo.GetTeacherByID(id)
}

// 教師情報を更新する
// 更新データを受け取り、データベースを更新して、更新後の教師情報を返す
func (s *TeacherService) Update(ctx context.Context, teacher *domain.Teacher) (*domain.Teacher, error) {
	// 教師情報を更新
	err := s.repo.UpdateTeacher(teacher)
	if err != nil {
		return nil, err
	}

	// 更新後の最新の教師情報を取得
	return s.repo.GetTeacherByID(teacher.ID)
}

// 指定されたIDの教師を削除する
func (s *TeacherService) Delete(ctx context.Context, id int64) error {
	return s.repo.DeleteTeacher(id)
}

// 教師のログイン認証を行う
// メールアドレスとパスワードを検証し、有効な場合はJWTトークンを返す
func (s *TeacherService) Login(ctx context.Context, email, password string) (string, error) {
	return s.repo.LoginTeacher(email, password)
}

// 教師に学生を割り当てる
func (s *TeacherService) AssignStudent(ctx context.Context, teacherID, studentID int64) error {
	return s.repo.AssignStudent(teacherID, studentID)
}

// 教師に割り当てられた学生一覧を取得する
func (s *TeacherService) GetStudents(ctx context.Context, teacherID int64) ([]domain.Student, error) {
	return s.repo.GetStudentsByTeacherID(teacherID)
}
