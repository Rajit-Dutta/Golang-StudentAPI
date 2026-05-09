package students

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Rajit-Dutta/StudentAPI/internal/types"
	"github.com/Rajit-Dutta/StudentAPI/internal/utils/responses"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
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

		responses.WriteJSON(w, http.StatusCreated, map[string]string{"success": "OK"})
	}
}
