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

// 学生リポジトリ構造体：データベースを使用した学生データの永続化を実装
type StudentRepository struct {
	// データベース接続
	db *gorm.DB
	// アプリケーション設定
	cfg *config.Config
}

// 学生データベースモデル：データベースのstudentsテーブルとマッピング
type Student struct {
	// 学生の一意識別子
	ID uint `gorm:"primaryKey"`
	// 学生の氏名
	Name string `gorm:"not null"`
	// 学生のメールアドレス（一意）
	Email string `gorm:"uniqueIndex;not null"`
	// 学生のパスワード（ハッシュ化）
	Password string `gorm:"not null"`
	// 学生の年齢
	Age int `gorm:"not null"`
	// レコードの作成日時
	CreatedAt time.Time `gorm:"autoCreateTime"`
	// レコードの更新日時
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	// 担当教師（多対多関係）
	Teachers []Teacher `gorm:"many2many:teacher_students;"`
}

// テーブル名を指定する
func (Student) TableName() string {
	return "students"
}

// 新しい学生リポジトリインスタンスを作成する
func NewStudentRepository(db *gorm.DB, cfg *config.Config) *StudentRepository {
	return &StudentRepository{
		db:  db,
		cfg: cfg,
	}
}

// 新しい学生を作成する
// パスワードをハッシュ化し、データベースに保存して、作成された学生のIDを返す
func (r *StudentRepository) CreateStudent(name, email string, age int, password string) (int64, error) {
	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	// 学生データを作成
	student := Student{
		Name:     name,
		Email:    email,
		Age:      age,
		Password: string(hashedPassword),
	}

	// データベースに保存
	result := r.db.Create(&student)
	if result.Error != nil {
		return 0, fmt.Errorf("failed to create student: %w", result.Error)
	}

	return int64(student.ID), nil
}

// 指定されたIDの学生を取得する
func (r *StudentRepository) GetStudentByID(id int64) (*domain.Student, error) {
	var student Student
	result := r.db.First(&student, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no student found with id: %d", id)
		}
		return nil, fmt.Errorf("query error: %w", result.Error)
	}

	// データベースモデルをドメインモデルに変換
	return &domain.Student{
		ID:        int64(student.ID),
		Name:      student.Name,
		Email:     student.Email,
		Age:       student.Age,
		CreatedAt: student.CreatedAt,
		UpdatedAt: student.UpdatedAt,
	}, nil
}

// 指定されたメールアドレスの学生を取得する
func (r *StudentRepository) GetStudentByEmail(email string) (*domain.Student, error) {
	var student Student
	result := r.db.Where("email = ?", email).First(&student)
	if result.Error != nil {
		return nil, result.Error
	}

	// データベースモデルをドメインモデルに変換
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

// 学生のログイン認証を行う
// メールアドレスとパスワードを検証し、有効な場合はJWTトークンを返す
func (r *StudentRepository) LoginStudent(email, password string) (string, error) {
	var student Student
	result := r.db.Where("email = ?", email).First(&student)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("no student found with email: %s", email)
		}
		return "", fmt.Errorf("query error: %w", result.Error)
	}

	// パスワードを検証
	if err := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	// JWTトークンを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   student.ID,
		"email": student.Email,
		"role":  "student",
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
