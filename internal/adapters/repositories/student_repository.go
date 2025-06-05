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

// Student represents the student entity in the database
type Student struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	Age       int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Teachers  []Teacher `gorm:"many2many:teacher_students;"`
}

func (Student) TableName() string {
	return "students"
}

func NewStudentRepository(db *gorm.DB, cfg *config.Config) *StudentRepository {
	return &StudentRepository{
		db:  db,
		cfg: cfg,
	}
}

func (r *StudentRepository) CreateStudent(name, email string, age int, password string) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	student := Student{
		Name:     name,
		Email:    email,
		Age:      age,
		Password: string(hashedPassword),
	}

	result := r.db.Create(&student)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to create student: %w", result.Error)
	}

	return int64(student.ID), nil
}

func (r *StudentRepository) GetStudentByID(id int64) (*domain.Student, error) {
	var student Student
	result := r.db.First(&student, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no student found with id: %d", id)
		}
		return nil, fmt.Errorf("query error: %w", result.Error)
	}

	return &domain.Student{
		ID:        int64(student.ID),
		Name:      student.Name,
		Email:     student.Email,
		Age:       student.Age,
		CreatedAt: student.CreatedAt,
		UpdatedAt: student.UpdatedAt,
	}, nil
}

func (r *StudentRepository) GetStudentByEmail(email string) (*domain.Student, error) {
	var student Student
	result := r.db.Where("email = ?", email).First(&student)
	if result.Error != nil {
		return nil, result.Error
	}

	return &domain.Student{
		ID:        int64(student.ID),
		Name:      student.Name,
		Email:     student.Email,
		Age:       student.Age,
		Password:  student.Password,
		CreatedAt: student.CreatedAt,
		UpdatedAt: student.UpdatedAt,
	}, nil
}

func (r *StudentRepository) LoginStudent(email, password string) (string, error) {
	var student Student
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
		"role":  "student",
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(r.cfg.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return tokenString, nil
}
