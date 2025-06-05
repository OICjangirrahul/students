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
		Name:     "John Smith",
		Email:    "john.smith@example.com",
		Password: "password123",
		Subject:  "Mathematics",
	}

	// Mock expectations
	mockRepo.On("CreateTeacher", teacher.Name, teacher.Email, teacher.Password, teacher.Subject).
		Return(int64(1), nil)

	expectedTeacher := &domain.Teacher{
		ID:        1,
		Name:      teacher.Name,
		Email:     teacher.Email,
		Subject:   teacher.Subject,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.On("GetTeacherByID", int64(1)).Return(expectedTeacher, nil)

	// Test
	result, err := service.Create(ctx, teacher)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedTeacher.ID, result.ID)
	assert.Equal(t, expectedTeacher.Name, result.Name)
	assert.Equal(t, expectedTeacher.Email, result.Email)
	assert.Equal(t, expectedTeacher.Subject, result.Subject)
	mockRepo.AssertExpectations(t)
}

func TestTeacherService_GetByID(t *testing.T) {
	// Setup
	mockRepo := new(mocks.TeacherRepository)
	service := NewTeacherService(mockRepo)
	ctx := context.Background()

	expectedTeacher := &domain.Teacher{
		ID:        1,
		Name:      "John Smith",
		Email:     "john.smith@example.com",
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
		Name:    "John Smith Updated",
		Email:   "john.updated@example.com",
		Subject: "Physics",
	}

	updatedTeacher := &domain.Teacher{
		ID:        1,
		Name:      teacher.Name,
		Email:     teacher.Email,
		Subject:   teacher.Subject,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock expectations
	mockRepo.On("UpdateTeacher", teacher).Return(nil)
	mockRepo.On("GetTeacherByID", int64(1)).Return(updatedTeacher, nil)

	// Test
	result, err := service.Update(ctx, teacher)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, teacher.ID, result.ID)
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

	email := "john.smith@example.com"
	password := "password123"
	expectedToken := "jwt-token"

	// Mock expectations
	mockRepo.On("LoginTeacher", email, password).Return(expectedToken, nil)

	// Test
	token, err := service.Login(ctx, email, password)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, expectedToken, token)
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

	teacherID := int64(1)
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
	mockRepo.On("GetStudentsByTeacherID", teacherID).Return(expectedStudents, nil)

	// Test
	students, err := service.GetStudents(ctx, teacherID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, students)
	assert.Equal(t, len(expectedStudents), len(students))
	for i, student := range students {
		assert.Equal(t, expectedStudents[i].ID, student.ID)
		assert.Equal(t, expectedStudents[i].Name, student.Name)
		assert.Equal(t, expectedStudents[i].Email, student.Email)
		assert.Equal(t, expectedStudents[i].Age, student.Age)
	}
	mockRepo.AssertExpectations(t)
}
