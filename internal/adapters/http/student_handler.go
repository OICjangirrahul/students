// Package http provides HTTP handlers for the API
package http

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/OICjangirrahul/students/internal/core/domain"
	"github.com/OICjangirrahul/students/internal/core/ports"
	"github.com/OICjangirrahul/students/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 学生ハンドラー構造体：学生に関するHTTPリクエストを処理
type StudentHandler struct {
	// 学生サービスインターフェース
	studentService ports.StudentService
}

// 新しい学生ハンドラーインスタンスを作成する
func NewStudentHandler(studentService ports.StudentService) *StudentHandler {
	return &StudentHandler{
		studentService: studentService,
	}
}

// 学生を作成する
// @Summary Create a new student
// @Description Create a new student with the provided information
// @Tags students
// @Accept json
// @Produce json
// @Param request body domain.Student true "Student information"
// @Success 201 {object} response.Response{data=domain.Student} "Created student"
// @Failure 400 {object} response.Response "Validation error"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/students [post]
func (h *StudentHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("creating a student")

		// リクエストボディから学生データを取得
		var student domain.Student
		if err := c.ShouldBindJSON(&student); err != nil {
			if errors.Is(err, io.EOF) {
				c.JSON(http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
				return
			}
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// バリデーション実行
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			c.JSON(http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		// 学生を作成
		createdStudent, err := h.studentService.Create(c.Request.Context(), &student)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("student created successfully", slog.String("studentId", fmt.Sprint(createdStudent.ID)))
		response.Success(c, http.StatusCreated, createdStudent)
	}
}

// 指定されたIDの学生を取得する
// @Summary Get a student by ID
// @Description Get a student's information by their ID
// @Tags students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Student ID"
// @Success 200 {object} response.Response{data=domain.Student} "Student found"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Student not found"
// @Router /api/v1/students/{id} [get]
func (h *StudentHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// パスパラメータからIDを取得
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		slog.Info("getting a student", slog.String("id", fmt.Sprint(id)))

		// 学生を取得
		student, err := h.studentService.GetByID(c.Request.Context(), id)
		if err != nil {
			slog.Error("error getting user", slog.String("id", fmt.Sprint(id)))
			c.JSON(http.StatusNotFound, response.GeneralError(err))
			return
		}

		response.Success(c, http.StatusOK, student)
	}
}

// 学生のログイン認証を行う
// @Summary Login student
// @Description Authenticate a student and return a JWT token
// @Tags students
// @Accept json
// @Produce json
// @Param request body domain.StudentLogin true "Login credentials"
// @Success 200 {object} response.Response{data=map[string]string{token=string}} "Login successful"
// @Failure 401 {object} response.Response "Invalid credentials"
// @Router /api/v1/students/login [post]
func (h *StudentHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("logging in a student")

		// リクエストボディからログイン情報を取得
		var login domain.StudentLogin
		if err := c.ShouldBindJSON(&login); err != nil {
			if errors.Is(err, io.EOF) {
				c.JSON(http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
				return
			}
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// バリデーション実行
		if err := validator.New().Struct(login); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			c.JSON(http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		// ログイン認証を実行
		token, err := h.studentService.Login(c.Request.Context(), login.Email, login.Password)
		if err != nil {
			slog.Error("error logging in", slog.String("email", login.Email), slog.String("error", err.Error()))
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid credentials")))
			return
		}

		slog.Info("student logged in successfully", slog.String("email", login.Email))
		response.Success(c, http.StatusOK, gin.H{"token": token})
	}
}
