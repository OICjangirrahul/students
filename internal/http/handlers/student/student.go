package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/OICjangirrahul/students/internal/storage"
	"github.com/OICjangirrahul/students/internal/types"
	"github.com/OICjangirrahul/students/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating a student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		//if body is empty
		if errors.Is(err, io.EOF) {
			// response.WriteJson(w,http.StatusBadRequest, response.GeneralError(err))

			//custom error message using fmt
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		//convert string into byte
		// w.Write([]byte("Welcome to students api"))

		if err := validator.New().Struct(student); err != nil {
			//type change err to => validateErrors
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
			student.Password,
		)

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
		}

		slog.Info("user created sucessfully", slog.String("userId", fmt.Sprint(lastId)))
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a student", slog.String("id", id))
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		}
		student, err := storage.GetStudentById(intId)

		if err != nil {
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}

// func Login(storage storage.Storage) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		slog.Info("logging in a student")

// 		var loginReq struct {
// 			Email    string `json:"email" validate:"required,email"`
// 			Password string `json:"password" validate:"required"`
// 		}

// 		err := json.NewDecoder(r.Body).Decode(&loginReq)
// 		if errors.Is(err, io.EOF) {
// 			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
// 			return
// 		}

// 		if err != nil {
// 			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
// 			return
// 		}

// 		if err := validator.New().Struct(loginReq); err != nil {
// 			validateErrs := err.(validator.ValidationErrors)
// 			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
// 			return
// 		}

// 		token, err := storage.LoginStudent(loginReq.Email, loginReq.Password)
// 		if err != nil {
// 			slog.Error("error logging in", slog.String("email", loginReq.Email), slog.String("error", err.Error()))
// 			response.WriteJson(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid credentials")))
// 			return
// 		}

// 		slog.Info("user logged in successfully", slog.String("email", loginReq.Email))
// 		response.WriteJson(w, http.StatusOK, map[string]string{"token": token})
// 	}
// }
