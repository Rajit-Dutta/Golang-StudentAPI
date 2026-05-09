package students

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Rajit-Dutta/StudentAPI/internal/types"
	"github.com/Rajit-Dutta/StudentAPI/internal/utils/responses"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			responses.WriteJSON(w, http.StatusBadGateway, err.Error())
			return
		}

		responses.WriteJSON(w, http.StatusCreated, map[string]string{"success": "OK"})

		slog.Info("creating a student")
	}
}
