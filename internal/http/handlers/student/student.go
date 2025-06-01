package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/OICjangirrahul/students/internal/types"
	"github.com/OICjangirrahul/students/internal/utils/response"
)

func New() http.HandlerFunc{
  	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating a student")
		
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		//if body is empty
		if errors.Is(err, io.EOF){
			// response.WriteJson(w,http.StatusBadRequest, response.GeneralError(err))

			//custom error message using fmt
			response.WriteJson(w,http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err !=  nil{
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		//convert string into byte
		// w.Write([]byte("Welcome to students api"))

		response.WriteJson(w,http.StatusCreated,map[string] string{"success": "OK"})
	}
}