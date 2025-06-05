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

// 教師ハンドラー構造体：教師に関するHTTPリクエストを処理
type TeacherHandler struct {
	// 教師サービスインターフェース
	teacherService ports.TeacherService
}

// 新しい教師ハンドラーインスタンスを作成する
func NewTeacherHandler(teacherService ports.TeacherService) *TeacherHandler {
	return &TeacherHandler{
		teacherService: teacherService,
	}
}

// 教師を作成する
// @Summary Create a new teacher
// @Description Create a new teacher with the provided information
// @Tags teachers
// @Accept json
// @Produce json
// @Param request body domain.Teacher true "Teacher information"
// @Success 201 {object} response.Response{data=domain.Teacher} "Created teacher"
// @Failure 400 {object} response.Response "Validation error"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /api/v1/teachers [post]
func (h *TeacherHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("creating a teacher")

		// リクエストボディから教師データを取得
		var teacher domain.Teacher
		if err := c.ShouldBindJSON(&teacher); err != nil {
			if errors.Is(err, io.EOF) {
				c.JSON(http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
				return
			}
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// バリデーション実行
		if err := validator.New().Struct(teacher); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			c.JSON(http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		// 教師を作成
		createdTeacher, err := h.teacherService.Create(c.Request.Context(), &teacher)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("teacher created successfully", slog.String("teacherId", fmt.Sprint(createdTeacher.ID)))
		response.Success(c, http.StatusCreated, createdTeacher)
	}
}

// 指定されたIDの教師を取得する
// @Summary Get a teacher by ID
// @Description Get a teacher's information by their ID
// @Tags teachers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Teacher ID"
// @Success 200 {object} response.Response{data=domain.Teacher} "Teacher found"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Teacher not found"
// @Router /api/v1/teachers/{id} [get]
func (h *TeacherHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// パスパラメータからIDを取得
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		slog.Info("getting a teacher", slog.String("id", fmt.Sprint(id)))

		// 教師を取得
		teacher, err := h.teacherService.GetByID(c.Request.Context(), id)
		if err != nil {
			slog.Error("error getting teacher", slog.String("id", fmt.Sprint(id)))
			c.JSON(http.StatusNotFound, response.GeneralError(err))
			return
		}

		response.Success(c, http.StatusOK, teacher)
	}
}

// 教師情報を更新する
// @Summary Update a teacher
// @Description Update a teacher's information
// @Tags teachers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Teacher ID"
// @Param request body domain.Teacher true "Teacher information"
// @Success 200 {object} response.Response{data=domain.Teacher} "Teacher updated"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Teacher not found"
// @Router /api/v1/teachers/{id} [put]
func (h *TeacherHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		// パスパラメータからIDを取得
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// リクエストボディから教師データを取得
		var teacher domain.Teacher
		if err := c.ShouldBindJSON(&teacher); err != nil {
			if errors.Is(err, io.EOF) {
				c.JSON(http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
				return
			}
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// IDを設定
		teacher.ID = id
		// 教師情報を更新
		updatedTeacher, err := h.teacherService.Update(c.Request.Context(), &teacher)
		if err != nil {
			slog.Error("error updating teacher", slog.String("id", fmt.Sprint(id)))
			c.JSON(http.StatusNotFound, response.GeneralError(err))
			return
		}

		response.Success(c, http.StatusOK, updatedTeacher)
	}
}

// 指定されたIDの教師を削除する
// @Summary Delete a teacher
// @Description Delete a teacher by their ID
// @Tags teachers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Teacher ID"
// @Success 204 "Teacher deleted"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Teacher not found"
// @Router /api/v1/teachers/{id} [delete]
func (h *TeacherHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		// パスパラメータからIDを取得
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// 教師を削除
		if err := h.teacherService.Delete(c.Request.Context(), id); err != nil {
			slog.Error("error deleting teacher", slog.String("id", fmt.Sprint(id)))
			c.JSON(http.StatusNotFound, response.GeneralError(err))
			return
		}

		response.Success(c, http.StatusNoContent, nil)
	}
}

// 教師のログイン認証を行う
// @Summary Login teacher
// @Description Authenticate a teacher and return a JWT token
// @Tags teachers
// @Accept json
// @Produce json
// @Param request body domain.TeacherLogin true "Login credentials"
// @Success 200 {object} response.Response{data=map[string]string{token=string}} "Login successful"
// @Failure 401 {object} response.Response "Invalid credentials"
// @Router /api/v1/teachers/login [post]
func (h *TeacherHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("logging in a teacher")

		// リクエストボディからログイン情報を取得
		var login domain.TeacherLogin
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
		token, err := h.teacherService.Login(c.Request.Context(), login.Email, login.Password)
		if err != nil {
			slog.Error("error logging in", slog.String("email", login.Email), slog.String("error", err.Error()))
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid credentials")))
			return
		}

		slog.Info("teacher logged in successfully", slog.String("email", login.Email))
		response.Success(c, http.StatusOK, gin.H{"token": token})
	}
}

// 教師に学生を割り当てる
// @Summary Assign student to teacher
// @Description Assign a student to a teacher
// @Tags teachers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Teacher ID"
// @Param studentId path int true "Student ID"
// @Success 200 {object} response.Response{data=map[string]string{message=string}} "Student assigned"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Teacher or student not found"
// @Router /api/v1/teachers/{id}/students/{studentId} [post]
func (h *TeacherHandler) AssignStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		// パスパラメータから教師IDと学生IDを取得
		teacherID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		studentID, err := strconv.ParseInt(c.Param("studentId"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// 学生を教師に割り当て
		if err := h.teacherService.AssignStudent(c.Request.Context(), teacherID, studentID); err != nil {
			slog.Error("error assigning student",
				slog.String("teacherId", fmt.Sprint(teacherID)),
				slog.String("studentId", fmt.Sprint(studentID)))
			c.JSON(http.StatusNotFound, response.GeneralError(err))
			return
		}

		response.Success(c, http.StatusOK, gin.H{"message": "Student assigned successfully"})
	}
}

// 教師に割り当てられた学生一覧を取得する
// @Summary Get teacher's students
// @Description Get all students assigned to a teacher
// @Tags teachers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Teacher ID"
// @Success 200 {object} response.Response{data=[]domain.Student} "List of students"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "Teacher not found"
// @Router /api/v1/teachers/{id}/students [get]
func (h *TeacherHandler) GetStudents() gin.HandlerFunc {
	return func(c *gin.Context) {
		// パスパラメータから教師IDを取得
		teacherID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// 教師に割り当てられた学生一覧を取得
		students, err := h.teacherService.GetStudents(c.Request.Context(), teacherID)
		if err != nil {
			slog.Error("error getting students", slog.String("teacherId", fmt.Sprint(teacherID)))
			c.JSON(http.StatusNotFound, response.GeneralError(err))
			return
		}

		response.Success(c, http.StatusOK, students)
	}
}
