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

type TeacherHandler struct {
	teacherService ports.TeacherService
}

func NewTeacherHandler(teacherService ports.TeacherService) *TeacherHandler {
	return &TeacherHandler{
		teacherService: teacherService,
	}
}

// Create handles teacher creation
// @Summary      Create a new teacher
// @Description  Create a new teacher with the provided information
// @Tags         teachers
// @Accept       json
// @Produce      json
// @Param        teacher body domain.Teacher true "Teacher information"
// @Success      201 {object} domain.Teacher
// @Router       /api/v1/teachers [post]
func (h *TeacherHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("creating a teacher")

		var teacher domain.Teacher
		if err := c.ShouldBindJSON(&teacher); err != nil {
			if errors.Is(err, io.EOF) {
				c.JSON(http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
				return
			}
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(teacher); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			c.JSON(http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastID, err := h.teacherService.Create(c.Request.Context(), &teacher)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("teacher created successfully", slog.String("teacherId", fmt.Sprint(lastID)))
		c.JSON(http.StatusCreated, gin.H{"id": lastID})
	}
}

// GetByID handles getting a teacher by ID
// @Summary      Get a teacher by ID
// @Description  Get a teacher's information by their ID
// @Tags         teachers
// @Accept       json
// @Produce      json
// @Param        id path int true "Teacher ID"
// @Success      200 {object} domain.Teacher
// @Router       /api/v1/teachers/{id} [get]
func (h *TeacherHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		teacher, err := h.teacherService.GetByID(c.Request.Context(), id)
		if err != nil {
			slog.Error("error getting teacher", slog.String("id", fmt.Sprint(id)))
			c.JSON(http.StatusNotFound, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, teacher)
	}
}

// Update handles updating a teacher
// @Summary      Update a teacher
// @Description  Update a teacher's information
// @Tags         teachers
// @Accept       json
// @Produce      json
// @Param        id path int true "Teacher ID"
// @Param        teacher body domain.Teacher true "Teacher information"
// @Success      200 {object} domain.Teacher
// @Router       /api/v1/teachers/{id} [put]
func (h *TeacherHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		var teacher domain.Teacher
		if err := c.ShouldBindJSON(&teacher); err != nil {
			if errors.Is(err, io.EOF) {
				c.JSON(http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
				return
			}
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		teacher.ID = id
		updatedTeacher, err := h.teacherService.Update(c.Request.Context(), &teacher)
		if err != nil {
			slog.Error("error updating teacher", slog.String("id", fmt.Sprint(id)))
			c.JSON(http.StatusNotFound, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, updatedTeacher)
	}
}

// Delete handles deleting a teacher
// @Summary      Delete a teacher
// @Description  Delete a teacher by their ID
// @Tags         teachers
// @Accept       json
// @Produce      json
// @Param        id path int true "Teacher ID"
// @Success      204 "No Content"
// @Router       /api/v1/teachers/{id} [delete]
func (h *TeacherHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := h.teacherService.Delete(c.Request.Context(), id); err != nil {
			slog.Error("error deleting teacher", slog.String("id", fmt.Sprint(id)))
			c.JSON(http.StatusNotFound, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusNoContent, nil)
	}
}

// Login handles teacher authentication
// @Summary      Login teacher
// @Description  Authenticate a teacher and return a JWT token
// @Tags         teachers
// @Accept       json
// @Produce      json
// @Param        credentials body domain.TeacherLogin true "Login credentials"
// @Router       /api/v1/teachers/login [post]
func (h *TeacherHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("logging in a teacher")

		var login domain.TeacherLogin
		if err := c.ShouldBindJSON(&login); err != nil {
			if errors.Is(err, io.EOF) {
				c.JSON(http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
				return
			}
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(login); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			c.JSON(http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		token, err := h.teacherService.Login(c.Request.Context(), login.Email, login.Password)
		if err != nil {
			slog.Error("error logging in", slog.String("email", login.Email), slog.String("error", err.Error()))
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid credentials")))
			return
		}

		slog.Info("teacher logged in successfully", slog.String("email", login.Email))
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

// AssignStudent handles assigning a student to a teacher
// @Summary      Assign student to teacher
// @Description  Assign a student to a teacher
// @Tags         teachers
// @Accept       json
// @Produce      json
// @Param        id path int true "Teacher ID"
// @Param        studentId path int true "Student ID"
// @Router       /api/v1/teachers/{id}/students/{studentId} [post]
func (h *TeacherHandler) AssignStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		if err := h.teacherService.AssignStudent(c.Request.Context(), teacherID, studentID); err != nil {
			slog.Error("error assigning student",
				slog.String("teacherId", fmt.Sprint(teacherID)),
				slog.String("studentId", fmt.Sprint(studentID)))
			c.JSON(http.StatusNotFound, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Student assigned successfully"})
	}
}

// GetStudents handles getting all students assigned to a teacher
// @Summary      Get teacher's students
// @Description  Get all students assigned to a teacher
// @Tags         teachers
// @Accept       json
// @Produce      json
// @Param        id path int true "Teacher ID"
// @Success      200 {array} domain.Student
// @Router       /api/v1/teachers/{id}/students [get]
func (h *TeacherHandler) GetStudents() gin.HandlerFunc {
	return func(c *gin.Context) {
		teacherID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		students, err := h.teacherService.GetStudents(c.Request.Context(), teacherID)
		if err != nil {
			slog.Error("error getting students", slog.String("teacherId", fmt.Sprint(teacherID)))
			c.JSON(http.StatusNotFound, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, students)
	}
}
