package ports

import "github.com/OICjangirrahul/students/internal/core/domain"

//go:generate mockery --name=StudentRepository --output=mocks --outpkg=mocks --case=snake
type StudentRepository interface {
	CreateStudent(name, email string, age int, password string) (int64, error)
	GetStudentByID(id int64) (*domain.Student, error)
	GetStudentByEmail(email string) (*domain.Student, error)
	LoginStudent(email, password string) (string, error)
}

//go:generate mockery --name=TeacherRepository --output=mocks --outpkg=mocks --case=snake
type TeacherRepository interface {
	CreateTeacher(name, email, password, subject string) (int64, error)
	GetTeacherByID(id int64) (*domain.Teacher, error)
	GetTeacherByEmail(email string) (*domain.Teacher, error)
	UpdateTeacher(teacher *domain.Teacher) error
	DeleteTeacher(id int64) error
	AssignStudent(teacherID, studentID int64) error
	GetStudentsByTeacherID(teacherID int64) ([]domain.Student, error)
	LoginTeacher(email, password string) (string, error)
}
