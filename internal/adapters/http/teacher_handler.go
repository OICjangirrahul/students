package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/OICjangirrahul/students/internal/core/domain"
	"github.com/OICjangirrahul/students/internal/core/services"
	"github.com/OICjangirrahul/students/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

type TeacherHandler struct {
	service *services.TeacherService
}

func NewTeacherHandler(service *services.TeacherService) *TeacherHandler {
	return &TeacherHandler{
		service: service,
	}
}

func (h *TeacherHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating a teacher")

		var teacher domain.Teacher
		err := json.NewDecoder(r.Body).Decode(&teacher)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(teacher); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastID, err := h.service.CreateTeacher(teacher)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("teacher created successfully", slog.String("teacherId", fmt.Sprint(lastID)))
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastID})
	}
}

func (h *TeacherHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a teacher", slog.String("id", id))

		intID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		teacher, err := h.service.GetTeacherByID(intID)
		if err != nil {
			slog.Error("error getting teacher", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, teacher)
	}
}

func (h *TeacherHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("updating a teacher", slog.String("id", id))

		intID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		var teacher domain.Teacher
		err = json.NewDecoder(r.Body).Decode(&teacher)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		err = h.service.UpdateTeacher(intID, teacher)
		if err != nil {
			slog.Error("error updating teacher", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"status": "success"})
	}
}

func (h *TeacherHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("deleting a teacher", slog.String("id", id))

		intID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		err = h.service.DeleteTeacher(intID)
		if err != nil {
			slog.Error("error deleting teacher", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"status": "success"})
	}
}

func (h *TeacherHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("logging in a teacher")

		var login domain.TeacherLogin
		err := json.NewDecoder(r.Body).Decode(&login)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(login); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		token, err := h.service.LoginTeacher(login)
		if err != nil {
			slog.Error("error logging in", slog.String("email", login.Email), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid credentials")))
			return
		}

		slog.Info("teacher logged in successfully", slog.String("email", login.Email))
		response.WriteJson(w, http.StatusOK, map[string]string{"token": token})
	}
}

func (h *TeacherHandler) AssignStudent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teacherID := r.PathValue("teacherId")
		studentID := r.PathValue("studentId")
		slog.Info("assigning student to teacher",
			slog.String("teacherId", teacherID),
			slog.String("studentId", studentID),
		)

		teacherIntID, err := strconv.ParseInt(teacherID, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid teacher ID")))
			return
		}

		studentIntID, err := strconv.ParseInt(studentID, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid student ID")))
			return
		}

		err = h.service.AssignStudentToTeacher(teacherIntID, studentIntID)
		if err != nil {
			slog.Error("error assigning student to teacher",
				slog.String("teacherId", teacherID),
				slog.String("studentId", studentID),
				slog.String("error", err.Error()),
			)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"status": "success"})
	}
}

func (h *TeacherHandler) GetStudents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teacherID := r.PathValue("teacherId")
		slog.Info("getting teacher's students", slog.String("teacherId", teacherID))

		teacherIntID, err := strconv.ParseInt(teacherID, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid teacher ID")))
			return
		}

		students, err := h.service.GetTeacherStudents(teacherIntID)
		if err != nil {
			slog.Error("error getting teacher's students",
				slog.String("teacherId", teacherID),
				slog.String("error", err.Error()),
			)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, students)
	}
}
