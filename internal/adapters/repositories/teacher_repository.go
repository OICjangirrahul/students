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

// 教師リポジトリ構造体：データベースを使用した教師データの永続化を実装
type TeacherRepository struct {
	// データベース接続
	db *gorm.DB
	// アプリケーション設定
	cfg *config.Config
}

// 教師データベースモデル：データベースのteachersテーブルとマッピング
type Teacher struct {
	// 教師の一意識別子
	ID uint `gorm:"primaryKey"`
	// 教師の氏名
	Name string `gorm:"not null"`
	// 教師のメールアドレス（一意）
	Email string `gorm:"uniqueIndex;not null"`
	// 教師のパスワード（ハッシュ化）
	Password string `gorm:"not null"`
	// 担当科目
	Subject string `gorm:"not null"`
	// レコードの作成日時
	CreatedAt time.Time `gorm:"autoCreateTime"`
	// レコードの更新日時
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	// 担当学生（多対多関係）
	Students []Student `gorm:"many2many:teacher_students;"`
}

// テーブル名を指定する
func (Teacher) TableName() string {
	return "teachers"
}

// 新しい教師リポジトリインスタンスを作成する
func NewTeacherRepository(db *gorm.DB, cfg *config.Config) *TeacherRepository {
	return &TeacherRepository{
		db:  db,
		cfg: cfg,
	}
}

// 新しい教師を作成する
// パスワードをハッシュ化し、データベースに保存して、作成された教師のIDを返す
func (r *TeacherRepository) CreateTeacher(name, email, password, subject string) (int64, error) {
	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	// 教師データを作成
	teacher := Teacher{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Subject:  subject,
	}

	// データベースに保存
	result := r.db.Create(&teacher)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to create teacher: %w", result.Error)
	}

	return int64(teacher.ID), nil
}

// 指定されたIDの教師を取得する
func (r *TeacherRepository) GetTeacherByID(id int64) (*domain.Teacher, error) {
	var teacher Teacher
	result := r.db.First(&teacher, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no teacher found with id: %d", id)
		}
		return nil, fmt.Errorf("query error: %w", result.Error)
	}

	// データベースモデルをドメインモデルに変換
	return &domain.Teacher{
		ID:        int64(teacher.ID),
		Name:      teacher.Name,
		Email:     teacher.Email,
		Subject:   teacher.Subject,
		CreatedAt: teacher.CreatedAt,
		UpdatedAt: teacher.UpdatedAt,
	}, nil
}

// 指定されたメールアドレスの教師を取得する
func (r *TeacherRepository) GetTeacherByEmail(email string) (*domain.Teacher, error) {
	var teacher Teacher
	result := r.db.Where("email = ?", email).First(&teacher)
	if result.Error != nil {
		return nil, result.Error
	}

	// データベースモデルをドメインモデルに変換
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

// 教師情報を更新する
func (r *TeacherRepository) UpdateTeacher(teacher *domain.Teacher) error {
	// 更新するモデルを作成
	model := Teacher{
		Name:    teacher.Name,
		Email:   teacher.Email,
		Subject: teacher.Subject,
	}
	model.ID = uint(teacher.ID)

	// データベースを更新
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

// 指定されたIDの教師を削除する
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

// 教師のログイン認証を行う
// メールアドレスとパスワードを検証し、有効な場合はJWTトークンを返す
func (r *TeacherRepository) LoginTeacher(email, password string) (string, error) {
	var teacher Teacher
	result := r.db.Where("email = ?", email).First(&teacher)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("no teacher found with email: %s", email)
		}
		return "", fmt.Errorf("query error: %w", result.Error)
	}

	// パスワードを検証
	if err := bcrypt.CompareHashAndPassword([]byte(teacher.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	// JWTトークンを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   teacher.ID,
		"email": teacher.Email,
		"role":  "teacher",
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	})

	// トークンに署名
	tokenString, err := token.SignedString([]byte(r.cfg.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return tokenString, nil
}

// 教師に学生を割り当てる
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

// 教師に割り当てられた学生一覧を取得する
func (r *TeacherRepository) GetStudentsByTeacherID(teacherID int64) ([]domain.Student, error) {
	var students []Student
	// 教師に割り当てられた学生を取得
	result := r.db.Raw(`
		SELECT s.* FROM students s
		JOIN teacher_students ts ON s.id = ts.student_id
		WHERE ts.teacher_id = ?
	`, teacherID).Scan(&students)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get students: %w", result.Error)
	}

	// データベースモデルをドメインモデルに変換
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
