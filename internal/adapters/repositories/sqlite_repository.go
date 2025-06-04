package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/OICjangirrahul/students/internal/config"
	"github.com/OICjangirrahul/students/internal/core/domain"
	"github.com/golang-jwt/jwt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type SQLiteRepository struct {
	db  *sql.DB
	cfg *config.Config
}

func NewSQLiteRepository(cfg *config.Config) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		age INTEGER,
		password TEXT
	)`)

	if err != nil {
		return nil, err
	}

	return &SQLiteRepository{
		db:  db,
		cfg: cfg,
	}, nil
}

func (r *SQLiteRepository) CreateStudent(name, email string, age int, password string) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	stmt, err := r.db.Prepare("INSERT INTO students (name, email, age, password) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age, string(hashedPassword))
	if err != nil {
		return 0, fmt.Errorf("failed to execute statement: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return lastID, nil
}

func (r *SQLiteRepository) GetStudentByID(id int64) (domain.Student, error) {
	stmt, err := r.db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return domain.Student{}, err
	}
	defer stmt.Close()

	var student domain.Student
	err = stmt.QueryRow(id).Scan(&student.ID, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Student{}, fmt.Errorf("no students found with id: %d", id)
		}
		return domain.Student{}, fmt.Errorf("query error: %w", err)
	}
	return student, nil
}

func (r *SQLiteRepository) LoginStudent(email, password string) (string, error) {
	if r.db == nil {
		return "", errors.New("database not initialized")
	}

	stmt, err := r.db.Prepare("SELECT id, password FROM students WHERE email = ? LIMIT 1")
	if err != nil {
		return "", fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var studentID int64
	var hashedPassword string
	err = stmt.QueryRow(email).Scan(&studentID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no student found with email: %s", email)
		}
		return "", fmt.Errorf("query error: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   studentID,
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(r.cfg.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return tokenString, nil
}
