package services

import (
	"context"
	"testing"
	"time"

	"github.com/OICjangirrahul/students/internal/core/domain"
	"github.com/OICjangirrahul/students/internal/core/ports/mocks"
	"github.com/stretchr/testify/assert"
)

func TestStudentService_Create(t *testing.T) {
	// Setup
	mockRepo := new(mocks.StudentRepository)
	service := NewStudentService(mockRepo)
	ctx := context.Background()

	student := &domain.Student{
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      20,
		Password: "password123",
	}

	// Mock expectations
	mockRepo.On("CreateStudent", student.Name, student.Email, student.Age, student.Password).
		Return(int64(1), nil)

	expectedStudent := &domain.Student{
		ID:        1,
		Name:      student.Name,
		Email:     student.Email,
		Age:       student.Age,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.On("GetStudentByID", int64(1)).Return(expectedStudent, nil)

	// Test
	result, err := service.Create(ctx, student)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedStudent.ID, result.ID)
	assert.Equal(t, expectedStudent.Name, result.Name)
	assert.Equal(t, expectedStudent.Email, result.Email)
	mockRepo.AssertExpectations(t)
}

func TestStudentService_GetByID(t *testing.T) {
	// Setup
	mockRepo := new(mocks.StudentRepository)
	service := NewStudentService(mockRepo)
	ctx := context.Background()

	expectedStudent := &domain.Student{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		Age:       20,
		CreatedAt: time.Now(),
	}

	// Mock expectations
	mockRepo.On("GetStudentByID", int64(1)).Return(expectedStudent, nil)

	// Test
	result, err := service.GetByID(ctx, 1)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedStudent.ID, result.ID)
	assert.Equal(t, expectedStudent.Name, result.Name)
	assert.Equal(t, expectedStudent.Email, result.Email)
	mockRepo.AssertExpectations(t)
}

func TestStudentService_Login(t *testing.T) {
	// Setup
	mockRepo := new(mocks.StudentRepository)
	service := NewStudentService(mockRepo)
	ctx := context.Background()

	email := "john@example.com"
	password := "password123"
	expectedToken := "jwt-token"

	// Mock expectations
	mockRepo.On("LoginStudent", email, password).Return(expectedToken, nil)

	// Test
	token, err := service.Login(ctx, email, password)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, expectedToken, token)
	mockRepo.AssertExpectations(t)
}
