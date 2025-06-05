package repositories

import (
	"fmt"
	"time"

	"github.com/OICjangirrahul/students/internal/config"
	"github.com/OICjangirrahul/students/internal/core/domain"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type StudentRepository struct {
	db  *gorm.DB
	cfg *config.Config
}

type StudentModel struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"not null"`
	Email     string         `gorm:"uniqueIndex;not null"`
	Age       int            `gorm:"not null"`
	Password  string         `gorm:"not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	Teachers  []TeacherModel `gorm:"many2many:teacher_students;"`
}

func NewStudentRepository(db *gorm.DB, cfg *config.Config) *StudentRepository {
	return &StudentRepository{
		db:  db,
		cfg: cfg,
	}
}

func (r *StudentRepository) CreateStudent(name, email string, age int, password string) (int64, error) {
	student := &domain.Student{
		Name:     name,
		Email:    email,
		Age:      age,
		Password: password,
	}

	result := r.db.Create(student)
	if result.Error != nil {
		return 0, result.Error
	}

	return student.ID, nil
}

func (r *StudentRepository) GetStudentByID(id int64) (*domain.Student, error) {
	var student domain.Student
	result := r.db.First(&student, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &student, nil
}

func (r *StudentRepository) GetStudentByEmail(email string) (*domain.Student, error) {
	var student domain.Student
	result := r.db.Where("email = ?", email).First(&student)
	if result.Error != nil {
		return nil, result.Error
	}
	return &student, nil
}

func (r *StudentRepository) LoginStudent(email, password string) (string, error) {
	var student StudentModel
	result := r.db.Where("email = ?", email).First(&student)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("no student found with email: %s", email)
		}
		return "", fmt.Errorf("query error: %w", result.Error)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   student.ID,
		"email": student.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(r.cfg.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return tokenString, nil
}
