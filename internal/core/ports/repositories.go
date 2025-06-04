package ports

import "github.com/OICjangirrahul/students/internal/core/domain"

type StudentRepository interface {
	CreateStudent(name string, email string, age int, password string) (int64, error)
	GetStudentByID(id int64) (domain.Student, error)
	LoginStudent(email, password string) (string, error)
}
