package services

import (
	"context"

	"github.com/OICjangirrahul/students/internal/core/domain"
	"github.com/OICjangirrahul/students/internal/core/ports"
)

type TeacherService struct {
	repo ports.TeacherRepository
}

func NewTeacherService(repo ports.TeacherRepository) *TeacherService {
	return &TeacherService{
		repo: repo,
	}
}

func (s *TeacherService) Create(ctx context.Context, teacher *domain.Teacher) (*domain.Teacher, error) {
	id, err := s.repo.CreateTeacher(teacher.Name, teacher.Email, teacher.Password, teacher.Subject)
	if err != nil {
		return nil, err
	}

	// Fetch the complete record to get timestamps
	return s.repo.GetTeacherByID(id)
}

func (s *TeacherService) GetByID(ctx context.Context, id int64) (*domain.Teacher, error) {
	return s.repo.GetTeacherByID(id)
}

func (s *TeacherService) Update(ctx context.Context, teacher *domain.Teacher) (*domain.Teacher, error) {
	err := s.repo.UpdateTeacher(teacher)
	if err != nil {
		return nil, err
	}

	// Fetch the updated record to get the latest timestamps
	return s.repo.GetTeacherByID(teacher.ID)
}

func (s *TeacherService) Delete(ctx context.Context, id int64) error {
	return s.repo.DeleteTeacher(id)
}

func (s *TeacherService) Login(ctx context.Context, email, password string) (string, error) {
	return s.repo.LoginTeacher(email, password)
}

func (s *TeacherService) AssignStudent(ctx context.Context, teacherID, studentID int64) error {
	return s.repo.AssignStudent(teacherID, studentID)
}

func (s *TeacherService) GetStudents(ctx context.Context, teacherID int64) ([]domain.Student, error) {
	return s.repo.GetStudentsByTeacherID(teacherID)
}
