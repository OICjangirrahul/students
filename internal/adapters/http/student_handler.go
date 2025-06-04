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

type StudentHandler struct {
	service *services.StudentService
}

func NewStudentHandler(service *services.StudentService) *StudentHandler {
	return &StudentHandler{
		service: service,
	}
}

func (h *StudentHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating a student")

		var student domain.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastID, err := h.service.CreateStudent(student)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("user created successfully", slog.String("userId", fmt.Sprint(lastID)))
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastID})
	}
}

func (h *StudentHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a student", slog.String("id", id))

		intID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := h.service.GetStudentByID(intID)
		if err != nil {
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}

func (h *StudentHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("logging in a student")

		var login domain.StudentLogin
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

		token, err := h.service.LoginStudent(login)
		if err != nil {
			slog.Error("error logging in", slog.String("email", login.Email), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid credentials")))
			return
		}

		slog.Info("user logged in successfully", slog.String("email", login.Email))
		response.WriteJson(w, http.StatusOK, map[string]string{"token": token})
	}
}
