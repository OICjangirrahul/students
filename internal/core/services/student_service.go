package services

import (
	"fmt"

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

func (s *StudentService) CreateStudent(student domain.Student) (int64, error) {
	return s.repo.CreateStudent(
		student.Name,
		student.Email,
		student.Age,
		student.Password,
	)
}

func (s *StudentService) GetStudentByID(id int64) (domain.Student, error) {
	student, err := s.repo.GetStudentByID(id)
	if err != nil {
		return domain.Student{}, fmt.Errorf("failed to get student: %w", err)
	}
	return student, nil
}

func (s *StudentService) LoginStudent(login domain.StudentLogin) (string, error) {
	token, err := s.repo.LoginStudent(login.Email, login.Password)
	if err != nil {
		return "", fmt.Errorf("failed to login: %w", err)
	}
	return token, nil
}
