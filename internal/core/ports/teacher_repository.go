package ports

import "github.com/OICjangirrahul/students/internal/core/domain"

type TeacherRepository interface {
	CreateTeacher(name, email, password, subject string) (int64, error)
	GetTeacherByID(id int64) (domain.Teacher, error)
	UpdateTeacher(id int64, name, subject string) error
	DeleteTeacher(id int64) error
	LoginTeacher(email, password string) (string, error)
	AssignStudentToTeacher(teacherID, studentID int64) error
	GetTeacherStudents(teacherID int64) ([]domain.Student, error)
}
