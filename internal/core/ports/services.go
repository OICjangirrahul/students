package ports

import (
	"context"

	"github.com/OICjangirrahul/students/internal/core/domain"
)

// StudentService defines the interface for student-related operations
type StudentService interface {
	Create(ctx context.Context, student *domain.Student) (*domain.Student, error)
	GetByID(ctx context.Context, id int64) (*domain.Student, error)
	Login(ctx context.Context, email, password string) (string, error)
}

// TeacherService defines the interface for teacher-related operations
type TeacherService interface {
	Create(ctx context.Context, teacher *domain.Teacher) (*domain.Teacher, error)
	GetByID(ctx context.Context, id int64) (*domain.Teacher, error)
	Update(ctx context.Context, teacher *domain.Teacher) (*domain.Teacher, error)
	Delete(ctx context.Context, id int64) error
	Login(ctx context.Context, email, password string) (string, error)
	AssignStudent(ctx context.Context, teacherID, studentID int64) error
	GetStudents(ctx context.Context, teacherID int64) ([]domain.Student, error)
}
