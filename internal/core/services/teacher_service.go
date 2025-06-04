package services

import (
	"fmt"

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

func (s *TeacherService) CreateTeacher(teacher domain.Teacher) (int64, error) {
	return s.repo.CreateTeacher(
		teacher.Name,
		teacher.Email,
		teacher.Password,
		teacher.Subject,
	)
}

func (s *TeacherService) GetTeacherByID(id int64) (domain.Teacher, error) {
	teacher, err := s.repo.GetTeacherByID(id)
	if err != nil {
		return domain.Teacher{}, fmt.Errorf("failed to get teacher: %w", err)
	}
	return teacher, nil
}

func (s *TeacherService) UpdateTeacher(id int64, teacher domain.Teacher) error {
	return s.repo.UpdateTeacher(id, teacher.Name, teacher.Subject)
}

func (s *TeacherService) DeleteTeacher(id int64) error {
	return s.repo.DeleteTeacher(id)
}

func (s *TeacherService) LoginTeacher(login domain.TeacherLogin) (string, error) {
	token, err := s.repo.LoginTeacher(login.Email, login.Password)
	if err != nil {
		return "", fmt.Errorf("failed to login: %w", err)
	}
	return token, nil
}

func (s *TeacherService) AssignStudentToTeacher(teacherID, studentID int64) error {
	return s.repo.AssignStudentToTeacher(teacherID, studentID)
}

func (s *TeacherService) GetTeacherStudents(teacherID int64) ([]domain.Student, error) {
	students, err := s.repo.GetTeacherStudents(teacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher's students: %w", err)
	}
	return students, nil
}
