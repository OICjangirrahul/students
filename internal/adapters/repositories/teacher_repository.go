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

type TeacherRepository struct {
	db  *gorm.DB
	cfg *config.Config
}
// TeacherModel represents the teacher entity in the database
type TeacherModel struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	Subject   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Students  []StudentModel `gorm:"many2many:teacher_students;"`
}

func NewTeacherRepository(db *gorm.DB, cfg *config.Config) *TeacherRepository {
	return &TeacherRepository{
		db:  db,
		cfg: cfg,
	}
}

func (r *TeacherRepository) CreateTeacher(name, email, password, subject string) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	teacher := TeacherModel{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Subject:  subject,
	}

	result := r.db.Create(&teacher)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to create teacher: %w", result.Error)
	}

	return int64(teacher.ID), nil
}

func (r *TeacherRepository) GetTeacherByID(id int64) (domain.Teacher, error) {
	var teacher TeacherModel
	result := r.db.Preload("Students").First(&teacher, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return domain.Teacher{}, fmt.Errorf("no teacher found with id: %d", id)
		}
		return domain.Teacher{}, fmt.Errorf("query error: %w", result.Error)
	}

	students := make([]domain.Student, len(teacher.Students))
	for i, s := range teacher.Students {
		students[i] = domain.Student{
			ID:    int64(s.ID),
			Name:  s.Name,
			Email: s.Email,
			Age:   s.Age,
		}
	}

	return domain.Teacher{
		ID:       int64(teacher.ID),
		Name:     teacher.Name,
		Email:    teacher.Email,
		Subject:  teacher.Subject,
		Students: students,
	}, nil
}

func (r *TeacherRepository) UpdateTeacher(id int64, name, subject string) error {
	result := r.db.Model(&TeacherModel{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":    name,
		"subject": subject,
	})

	if result.Error != nil {
		return fmt.Errorf("failed to update teacher: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no teacher found with id: %d", id)
	}

	return nil
}

func (r *TeacherRepository) DeleteTeacher(id int64) error {
	result := r.db.Delete(&TeacherModel{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete teacher: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no teacher found with id: %d", id)
	}

	return nil
}

func (r *TeacherRepository) LoginTeacher(email, password string) (string, error) {
	var teacher TeacherModel
	result := r.db.Where("email = ?", email).First(&teacher)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("no teacher found with email: %s", email)
		}
		return "", fmt.Errorf("query error: %w", result.Error)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(teacher.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   teacher.ID,
		"email": teacher.Email,
		"role":  "teacher",
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(r.cfg.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return tokenString, nil
}

func (r *TeacherRepository) AssignStudentToTeacher(teacherID, studentID int64) error {
	result := r.db.Exec("INSERT INTO teacher_students (teacher_id, student_id) VALUES (?, ?)", teacherID, studentID)
	if result.Error != nil {
		return fmt.Errorf("failed to assign student to teacher: %w", result.Error)
	}

	return nil
}

func (r *TeacherRepository) GetTeacherStudents(teacherID int64) ([]domain.Student, error) {
	var teacher TeacherModel
	result := r.db.Preload("Students").First(&teacher, teacherID)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get teacher: %w", result.Error)
	}

	students := make([]domain.Student, len(teacher.Students))
	for i, s := range teacher.Students {
		students[i] = domain.Student{
			ID:    int64(s.ID),
			Name:  s.Name,
			Email: s.Email,
			Age:   s.Age,
		}
	}

	return students, nil
}
