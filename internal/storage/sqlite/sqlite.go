package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/OICjangirrahul/students/internal/config"
	"github.com/OICjangirrahul/students/internal/types"
	"github.com/golang-jwt/jwt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Sqlite struct {
	Db  *sql.DB
	Cfg *config.Config
}

func NewSqlite(cfg *config.Config) (*Sqlite, error) {
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
	return &Sqlite{
		Db: db,
		Cfg: cfg,
	}, nil

}

func (s *Sqlite) CreateStudent(name, email string, age int, password string) (int64, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age, password) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age, string(hashedPassword))
	if err != nil {
		return 0, fmt.Errorf("failed to execute statement: %w", err)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return lastId, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no students found with id: %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error: %s", err)
	}
	return student, err

}
func (s *Sqlite) LoginStudent(email, password string) (string, error) {
    // Check for nil database
    if s == nil || s.Db == nil {
        return "", errors.New("database not initialized")
    }

    // Prepare SQL statement
    stmt, err := s.Db.Prepare("SELECT id, password FROM students WHERE email = ? LIMIT 1")
    if err != nil {
        return "", fmt.Errorf("failed to prepare statement: %w", err)
    }

    defer stmt.Close()

    // Query and scan results
    var studentId int64
    var hashedPassword string
		fmt.Println("e......")
    err = stmt.QueryRow(email).Scan(&studentId, &hashedPassword)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("no student found with email: %s", email)
        }
        return "", fmt.Errorf("query error: %w", err)
    }
    // Verify password with bcrypt
    if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
        return "", fmt.Errorf("invalid credentials: %w", err)
    }
    // Generate JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub":   studentId,
        "email": email,
        "exp":   time.Now().Add(24 * time.Hour).Unix(),
        "iat":   time.Now().Unix(),
    })

    tokenString, err := token.SignedString([]byte(s.Cfg.JWT.Secret))
    if err != nil {
        return "", fmt.Errorf("failed to generate JWT: %w", err)
    }

    return tokenString, nil
}