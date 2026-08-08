package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/JayantRautela/students-api/internal/storage"
	"github.com/JayantRautela/students-api/internal/types"
	"github.com/JayantRautela/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student")
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// request validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadGateway, response.ValidationError(validateErrs))
			return
		}

		id, err := storage.CreateStudent(student.Name, student.Email, student.Age)

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		slog.Info("User create successfully", slog.String("userId", fmt.Sprint(id)))

		response.WriteJson(w, http.StatusCreated, map[string]int{"id": id})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Fteching student details", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 0)

		if err != nil {
			slog.Error("Error parsing id")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return 
		}

		student, err := storage.GetStudentById(int(intId))

		if err != nil {
			slog.Error("Error getting student", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return 
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting all students")

		students, err := storage.GetStudents()

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return 
		}

		response.WriteJson(w, http.StatusOK, students)
	}
}

func DeleteStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request)  {
		slog.Info("Deleting Student Detail")
		id := r.PathValue("id")

		intId, err := strconv.ParseInt(id, 10, 0)

		if err != nil {
			slog.Error("Error parsing id")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return 
		}

		err = storage.DeleteStudent(int(intId))

		if err != nil {
			slog.Error("Error deleting student", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return 
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"message": "student details deleted"})
	}
}
