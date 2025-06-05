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

// Teacher represents the teacher entity in the database
type Teacher struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	Subject   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Students  []Student `gorm:"many2many:teacher_students;"`
}

func (Teacher) TableName() string {
	return "teachers"
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

	teacher := Teacher{
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

func (r *TeacherRepository) GetTeacherByID(id int64) (*domain.Teacher, error) {
	var teacher Teacher
	result := r.db.First(&teacher, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no teacher found with id: %d", id)
		}
		return nil, fmt.Errorf("query error: %w", result.Error)
	}

	return &domain.Teacher{
		ID:        int64(teacher.ID),
		Name:      teacher.Name,
		Email:     teacher.Email,
		Subject:   teacher.Subject,
		CreatedAt: teacher.CreatedAt,
		UpdatedAt: teacher.UpdatedAt,
	}, nil
}

func (r *TeacherRepository) GetTeacherByEmail(email string) (*domain.Teacher, error) {
	var teacher Teacher
	result := r.db.Where("email = ?", email).First(&teacher)
	if result.Error != nil {
		return nil, result.Error
	}

	return &domain.Teacher{
		ID:        int64(teacher.ID),
		Name:      teacher.Name,
		Email:     teacher.Email,
		Subject:   teacher.Subject,
		Password:  teacher.Password,
		CreatedAt: teacher.CreatedAt,
		UpdatedAt: teacher.UpdatedAt,
	}, nil
}

func (r *TeacherRepository) UpdateTeacher(teacher *domain.Teacher) error {
	model := Teacher{
		Name:    teacher.Name,
		Email:   teacher.Email,
		Subject: teacher.Subject,
	}
	model.ID = uint(teacher.ID)

	result := r.db.Model(&Teacher{}).Where("id = ?", teacher.ID).Updates(map[string]interface{}{
		"name":    teacher.Name,
		"email":   teacher.Email,
		"subject": teacher.Subject,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to update teacher: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no teacher found with id: %d", teacher.ID)
	}

	return nil
}

func (r *TeacherRepository) DeleteTeacher(id int64) error {
	result := r.db.Delete(&Teacher{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete teacher: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no teacher found with id: %d", id)
	}

	return nil
}

func (r *TeacherRepository) LoginTeacher(email, password string) (string, error) {
	var teacher Teacher
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

func (r *TeacherRepository) AssignStudent(teacherID, studentID int64) error {
	result := r.db.Exec(`
		INSERT INTO teacher_students (teacher_id, student_id) 
		VALUES (?, ?) 
		ON CONFLICT (teacher_id, student_id) DO NOTHING
	`, teacherID, studentID)
	if result.Error != nil {
		return fmt.Errorf("failed to assign student to teacher: %w", result.Error)
	}

	return nil
}

func (r *TeacherRepository) GetStudentsByTeacherID(teacherID int64) ([]domain.Student, error) {
	var students []Student
	result := r.db.Raw(`
		SELECT s.* FROM students s
		JOIN teacher_students ts ON s.id = ts.student_id
		WHERE ts.teacher_id = ?
	`, teacherID).Scan(&students)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get students: %w", result.Error)
	}

	// Convert to domain models
	domainStudents := make([]domain.Student, len(students))
	for i, s := range students {
		domainStudents[i] = domain.Student{
			ID:        int64(s.ID),
			Name:      s.Name,
			Email:     s.Email,
			Age:       s.Age,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		}
	}

	return domainStudents, nil
}
