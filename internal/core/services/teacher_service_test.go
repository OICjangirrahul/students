package services

import (
	"context"
	"testing"
	"time"

	"github.com/OICjangirrahul/students/internal/core/domain"
	"github.com/OICjangirrahul/students/internal/core/ports/mocks"
	"github.com/stretchr/testify/assert"
)

func TestTeacherService_Create(t *testing.T) {
	// Setup
	mockRepo := new(mocks.TeacherRepository)
	service := NewTeacherService(mockRepo)
	ctx := context.Background()

	teacher := &domain.Teacher{
		Name:     "Jane Smith",
		Email:    "jane@example.com",
		Subject:  "Mathematics",
		Password: "password123",
	}

	// Mock expectations
	mockRepo.On("CreateTeacher", teacher.Name, teacher.Email, teacher.Password, teacher.Subject).
		Return(int64(1), nil)

	// Test
	result, err := service.Create(ctx, teacher)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	mockRepo.AssertExpectations(t)
}

func TestTeacherService_GetByID(t *testing.T) {
	// Setup
	mockRepo := new(mocks.TeacherRepository)
	service := NewTeacherService(mockRepo)
	ctx := context.Background()

	expectedTeacher := &domain.Teacher{
		ID:        1,
		Name:      "Jane Smith",
		Email:     "jane@example.com",
		Subject:   "Mathematics",
		CreatedAt: time.Now(),
	}

	// Mock expectations
	mockRepo.On("GetTeacherByID", int64(1)).Return(expectedTeacher, nil)

	// Test
	result, err := service.GetByID(ctx, 1)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTeacher.ID, result.ID)
	assert.Equal(t, expectedTeacher.Name, result.Name)
	assert.Equal(t, expectedTeacher.Email, result.Email)
	assert.Equal(t, expectedTeacher.Subject, result.Subject)
	mockRepo.AssertExpectations(t)
}

func TestTeacherService_Update(t *testing.T) {
	// Setup
	mockRepo := new(mocks.TeacherRepository)
	service := NewTeacherService(mockRepo)
	ctx := context.Background()

	teacher := &domain.Teacher{
		ID:      1,
		Name:    "Jane Smith Updated",
		Email:   "jane.updated@example.com",
		Subject: "Physics",
	}

	// Mock expectations
	mockRepo.On("UpdateTeacher", teacher).Return(nil)

	// Test
	result, err := service.Update(ctx, teacher)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, teacher.Name, result.Name)
	assert.Equal(t, teacher.Email, result.Email)
	assert.Equal(t, teacher.Subject, result.Subject)
	mockRepo.AssertExpectations(t)
}

func TestTeacherService_Delete(t *testing.T) {
	// Setup
	mockRepo := new(mocks.TeacherRepository)
	service := NewTeacherService(mockRepo)
	ctx := context.Background()

	// Mock expectations
	mockRepo.On("DeleteTeacher", int64(1)).Return(nil)

	// Test
	err := service.Delete(ctx, 1)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTeacherService_Login(t *testing.T) {
	// Setup
	mockRepo := new(mocks.TeacherRepository)
	service := NewTeacherService(mockRepo)
	ctx := context.Background()

	email := "jane@example.com"
	password := "password123"
	expectedToken := "jwt-token"

	// Mock expectations
	mockRepo.On("GetTeacherByEmail", email).Return(&domain.Teacher{
		ID:       1,
		Email:    email,
		Password: password,
	}, nil)

	// Test
	token, err := service.Login(ctx, email, password)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotEqual(t, expectedToken, token) // Token should be dynamically generated
	mockRepo.AssertExpectations(t)
}

func TestTeacherService_AssignStudent(t *testing.T) {
	// Setup
	mockRepo := new(mocks.TeacherRepository)
	service := NewTeacherService(mockRepo)
	ctx := context.Background()

	teacherID := int64(1)
	studentID := int64(2)

	// Mock expectations
	mockRepo.On("AssignStudent", teacherID, studentID).Return(nil)

	// Test
	err := service.AssignStudent(ctx, teacherID, studentID)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTeacherService_GetStudents(t *testing.T) {
	// Setup
	mockRepo := new(mocks.TeacherRepository)
	service := NewTeacherService(mockRepo)
	ctx := context.Background()

	expectedStudents := []domain.Student{
		{
			ID:    1,
			Name:  "Student 1",
			Email: "student1@example.com",
			Age:   20,
		},
		{
			ID:    2,
			Name:  "Student 2",
			Email: "student2@example.com",
			Age:   21,
		},
	}

	// Mock expectations
	mockRepo.On("GetStudentsByTeacherID", int64(1)).Return(expectedStudents, nil)

	// Test
	result, err := service.GetStudents(ctx, 1)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(expectedStudents), len(result))
	assert.Equal(t, expectedStudents[0].ID, result[0].ID)
	assert.Equal(t, expectedStudents[1].ID, result[1].ID)
	mockRepo.AssertExpectations(t)
}
