package ports

import "github.com/OICjangirrahul/students/internal/core/domain"

// 学生リポジトリインターフェース：学生データの永続化操作を定義
//
//go:generate mockery --name=StudentRepository --output=mocks --outpkg=mocks --case=snake
type StudentRepository interface {
	// 新しい学生を作成し、作成された学生のIDを返す
	CreateStudent(name, email string, age int, password string) (int64, error)
	// 指定されたIDの学生を取得する
	GetStudentByID(id int64) (*domain.Student, error)
	// 指定されたメールアドレスの学生を取得する
	GetStudentByEmail(email string) (*domain.Student, error)
	// 学生のログイン認証を行い、JWTトークンを返す
	LoginStudent(email, password string) (string, error)
}

// 教師リポジトリインターフェース：教師データの永続化操作を定義
//
//go:generate mockery --name=TeacherRepository --output=mocks --outpkg=mocks --case=snake
type TeacherRepository interface {
	// 新しい教師を作成し、作成された教師のIDを返す
	CreateTeacher(name, email, password, subject string) (int64, error)
	// 指定されたIDの教師を取得する
	GetTeacherByID(id int64) (*domain.Teacher, error)
	// 指定されたメールアドレスの教師を取得する
	GetTeacherByEmail(email string) (*domain.Teacher, error)
	// 教師情報を更新する
	UpdateTeacher(teacher *domain.Teacher) error
	// 指定されたIDの教師を削除する
	DeleteTeacher(id int64) error
	// 教師に学生を割り当てる
	AssignStudent(teacherID, studentID int64) error
	// 教師に割り当てられた学生一覧を取得する
	GetStudentsByTeacherID(teacherID int64) ([]domain.Student, error)
	// 教師のログイン認証を行い、JWTトークンを返す
	LoginTeacher(email, password string) (string, error)
}
