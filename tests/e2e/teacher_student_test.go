package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/OICjangirrahul/students/internal"
	"github.com/OICjangirrahul/students/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	baseURL     string
	teacherData map[string]interface{}
	studentData map[string]interface{}
	authToken   string
	server      *http.Server
)

func setupTestServer(t *testing.T) (*gin.Engine, error) {
	// Load test configuration
	cfg, err := config.LoadConfig("../../config/test.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	handlers, err := internal.InitializeAppHandlers(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize handlers: %w", err)
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

	// Student routes
	students := v1.Group("/students")
	{
		students.POST("", handlers.Student.Create())
	}

	return r, nil
}

func TestMain(m *testing.M) {
	// Setup
	gin.SetMode(gin.TestMode)
	r, err := setupTestServer(nil)
	if err != nil {
		fmt.Printf("Failed to setup test server: %v\n", err)
		os.Exit(1)
	}

	// Start the server
	server = &http.Server{
		Addr:    ":8081",
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Failed to start server: %v\n", err)
			os.Exit(1)
		}
	}()

	// Wait for server to start
	time.Sleep(2 * time.Second)
	baseURL = "http://localhost:8081/api/v1"

	// Run tests
	code := m.Run()

	// Cleanup
	if err := server.Close(); err != nil {
		fmt.Printf("Failed to close server: %v\n", err)
	}

	os.Exit(code)
}

func TestE2E_TeacherStudentFlow(t *testing.T) {
	t.Run("1. Create Teacher", func(t *testing.T) {
		payload := map[string]string{
			"name":     "Jane Smith",
			"email":    "jane@example.com",
			"subject":  "Mathematics",
			"password": "password123",
		}

		resp, err := makeRequest("POST", "/teachers", payload, "")
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		teacherData = response["data"].(map[string]interface{})
		assert.NotEmpty(t, teacherData["id"])
	})

	t.Run("2. Teacher Login", func(t *testing.T) {
		payload := map[string]string{
			"email":    "jane@example.com",
			"password": "password123",
		}

		resp, err := makeRequest("POST", "/teachers/login", payload, "")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		data := response["data"].(map[string]interface{})
		authToken = data["token"].(string)
		assert.NotEmpty(t, authToken)
	})

	t.Run("3. Create Student", func(t *testing.T) {
		payload := map[string]string{
			"name":     "John Doe",
			"email":    "john@example.com",
			"password": "password123",
			"age":      "20",
		}

		resp, err := makeRequest("POST", "/students", payload, authToken)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		studentData = response["data"].(map[string]interface{})
		assert.NotEmpty(t, studentData["id"])
	})

	t.Run("4. Assign Student to Teacher", func(t *testing.T) {
		teacherID := fmt.Sprint(teacherData["id"])
		studentID := fmt.Sprint(studentData["id"])
		path := fmt.Sprintf("/teachers/%s/students/%s", teacherID, studentID)

		resp, err := makeRequest("POST", path, nil, authToken)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("5. Get Teacher's Students", func(t *testing.T) {
		teacherID := fmt.Sprint(teacherData["id"])
		path := fmt.Sprintf("/teachers/%s/students", teacherID)

		resp, err := makeRequest("GET", path, nil, authToken)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		students := response["data"].([]interface{})
		assert.NotEmpty(t, students)
		assert.Equal(t, 1, len(students))

		student := students[0].(map[string]interface{})
		assert.Equal(t, fmt.Sprint(studentData["id"]), fmt.Sprint(student["id"]))
	})

	t.Run("6. Update Teacher", func(t *testing.T) {
		teacherID := fmt.Sprint(teacherData["id"])
		path := fmt.Sprintf("/teachers/%s", teacherID)

		payload := map[string]string{
			"name":    "Jane Smith Updated",
			"email":   "jane.updated@example.com",
			"subject": "Physics",
		}

		resp, err := makeRequest("PUT", path, payload, authToken)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		data := response["data"].(map[string]interface{})
		assert.Equal(t, payload["name"], data["name"])
		assert.Equal(t, payload["email"], data["email"])
		assert.Equal(t, payload["subject"], data["subject"])
	})

	t.Run("7. Delete Teacher", func(t *testing.T) {
		teacherID := fmt.Sprint(teacherData["id"])
		path := fmt.Sprintf("/teachers/%s", teacherID)

		resp, err := makeRequest("DELETE", path, nil, authToken)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

func makeRequest(method, path string, payload interface{}, token string) (*http.Response, error) {
	var body []byte
	var err error

	if payload != nil {
		body, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, baseURL+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	return client.Do(req)
}
