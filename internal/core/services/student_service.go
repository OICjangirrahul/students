package services

import (
	"context"

	"github.com/OICjangirrahul/students/internal/core/domain"
	"github.com/OICjangirrahul/students/internal/core/ports"
)

type StudentService struct {
	repo ports.StudentRepository
}

func NewStudentService(repo ports.StudentRepository) *StudentService {
	return &StudentService{
		repo: repo,
	}
}

func (s *StudentService) Create(ctx context.Context, student *domain.Student) (*domain.Student, error) {
	id, err := s.repo.CreateStudent(student.Name, student.Email, student.Age, student.Password)
	if err != nil {
		return nil, err
	}

	student.ID = id
	return student, nil
}

func (s *StudentService) GetByID(ctx context.Context, id int64) (*domain.Student, error) {
	return s.repo.GetStudentByID(id)
}

func (s *StudentService) Login(ctx context.Context, email, password string) (string, error) {
	student, err := s.repo.GetStudentByEmail(email)
	if err != nil {
		return "", err
	}

	// TODO: Add password verification
	if student.Password != password {
		return "", domain.ErrInvalidCredentials
	}

	// TODO: Generate JWT token
	return "student-jwt-token", nil
}
