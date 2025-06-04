package storage

import "github.com/OICjangirrahul/students/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int, password string) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	LoginStudent(email, password string) (string, error)
}
