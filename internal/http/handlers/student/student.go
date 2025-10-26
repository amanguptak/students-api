package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/amanguptak/students-api/internal/types"
	"github.com/amanguptak/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		if err != nil {
			var ute *json.UnmarshalTypeError

			// Unmarshal Type Error

			// When JSON decoding fails because a field’s type doesn’t match (like a string instead of an int), Go returns an error of this type.
			if errors.As(err, &ute) {
				// e.g. "Age should be int, not string"
				msg := fmt.Sprintf("%s should be %s, not %s", ute.Field, ute.Type.String(), ute.Value)
				response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("%s", msg)))
				return
			}

			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

	

		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors) // doing typecasting here
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		// w.Write([]byte("Welcome to student api end point"))

		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "OK"})
	}
}
