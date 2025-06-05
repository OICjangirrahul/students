package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OICjangirrahul/students/internal"
	"github.com/OICjangirrahul/students/internal/config"
	"github.com/OICjangirrahul/students/internal/core/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T) (*gin.Engine, error) {
	// Load test configuration
	cfg, err := config.LoadConfig("../../config/test.yaml")
	if err != nil {
		return nil, err
	}

	handlers, err := internal.InitializeAppHandlers(cfg)
	if err != nil {
		return nil, err
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")

	// Teacher routes
	teachers := v1.Group("/teachers")
	{
		teachers.POST("", handlers.Teacher.Create())
		teachers.POST("/login", handlers.Teacher.Login())

		teacherManagement := teachers.Group("/:id")
		{
			teacherManagement.GET("", handlers.Teacher.GetByID())
			teacherManagement.PUT("", handlers.Teacher.Update())
			teacherManagement.DELETE("", handlers.Teacher.Delete())

			studentManagement := teacherManagement.Group("/students")
			{
				studentManagement.GET("", handlers.Teacher.GetStudents())
				studentManagement.POST("/:studentId", handlers.Teacher.AssignStudent())
			}
		}
	}

	return r, nil
}

func TestTeacherAPI_Create(t *testing.T) {
	r, err := setupTestServer(t)
	require.NoError(t, err)

	teacher := domain.Teacher{
		Name:     "Jane Smith",
		Email:    "jane@example.com",
		Subject:  "Mathematics",
		Password: "password123",
	}

	body, err := json.Marshal(teacher)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/teachers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Data struct {
			ID      int64  `json:"id"`
			Name    string `json:"name"`
			Email   string `json:"email"`
			Subject string `json:"subject"`
		} `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, teacher.Name, response.Data.Name)
	assert.Equal(t, teacher.Email, response.Data.Email)
	assert.Equal(t, teacher.Subject, response.Data.Subject)
}

func TestTeacherAPI_Login(t *testing.T) {
	r, err := setupTestServer(t)
	require.NoError(t, err)

	loginReq := domain.TeacherLogin{
		Email:    "jane@example.com",
		Password: "password123",
	}

	body, err := json.Marshal(loginReq)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/teachers/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.NotEmpty(t, response.Data.Token)
}

func TestTeacherAPI_GetByID(t *testing.T) {
	r, err := setupTestServer(t)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/teachers/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data struct {
			ID      int64  `json:"id"`
			Name    string `json:"name"`
			Email   string `json:"email"`
			Subject string `json:"subject"`
		} `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.NotEmpty(t, response.Data.ID)
	assert.NotEmpty(t, response.Data.Name)
	assert.NotEmpty(t, response.Data.Email)
}

func TestTeacherAPI_Update(t *testing.T) {
	r, err := setupTestServer(t)
	require.NoError(t, err)

	updateReq := map[string]string{
		"name":    "Jane Smith Updated",
		"email":   "jane.updated@example.com",
		"subject": "Physics",
	}

	body, err := json.Marshal(updateReq)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/teachers/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data struct {
			Name    string `json:"name"`
			Email   string `json:"email"`
			Subject string `json:"subject"`
		} `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, updateReq["name"], response.Data.Name)
	assert.Equal(t, updateReq["email"], response.Data.Email)
	assert.Equal(t, updateReq["subject"], response.Data.Subject)
}

func TestTeacherAPI_Delete(t *testing.T) {
	r, err := setupTestServer(t)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/teachers/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestTeacherAPI_AssignStudent(t *testing.T) {
	r, err := setupTestServer(t)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/teachers/1/students/2", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTeacherAPI_GetStudents(t *testing.T) {
	r, err := setupTestServer(t)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/teachers/1/students", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Data []struct {
			ID    int64  `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
			Age   int    `json:"age"`
		} `json:"data"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.NotEmpty(t, response.Data)
}
