package students

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Rajit-Dutta/StudentAPI/internal/storage"
	"github.com/Rajit-Dutta/StudentAPI/internal/types"
	"github.com/Rajit-Dutta/StudentAPI/internal/utils/responses"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating a student")
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			responses.WriteJSON(w, http.StatusBadRequest, responses.GeneralError(err))
			return
		}

		if err := validator.New().Struct(student); err != nil {
			validatorErrors := err.(validator.ValidationErrors)
			responses.WriteJSON(w, http.StatusBadRequest, responses.ValidatorError(validatorErrors))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Age,
			student.Email,
		)
		if err != nil {
			responses.WriteJSON(w, http.StatusInternalServerError, responses.GeneralError(err))
			return
		}

		slog.Info("Create a strudent entry", slog.String("userId:", fmt.Sprint(lastId)))

		responses.WriteJSON(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetByID(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a student", slog.String("id:", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			responses.WriteJSON(w, http.StatusBadRequest, responses.GeneralError(err))
			return
		}

		student, err := storage.GetStudentbyID(intId)
		if err != nil {
			responses.WriteJSON(w, http.StatusInternalServerError, responses.GeneralError(err))
		}

		responses.WriteJSON(w, http.StatusOK, student)
	}
}

func GetStudentsList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		students, err := storage.GetStudentList()
		if err != nil {
			responses.WriteJSON(w, http.StatusInternalServerError, responses.GeneralError(err))
			return
		}
		responses.WriteJSON(w, http.StatusOK, students)
	}
}
