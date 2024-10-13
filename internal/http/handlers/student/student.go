package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/PratikPradhan987/learn-go/internal/storage"
	"github.com/PratikPradhan987/learn-go/internal/types"
	"github.com/PratikPradhan987/learn-go/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		slog.Info("Creating a new student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		
		if errors.Is(err, io.EOF) {
			// response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
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
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
            return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
            student.Email,
            student.Age,
            // student.Grade,
            // student.Course,
            // student.Address,
            // student.PhoneNumber,
		)
		
		slog.Info("User Created Successfully", slog.String("userId", fmt.Sprint(lastId))) 

        if err!= nil {
            response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
            return
        }

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		id := r.PathValue("id")
		slog.Info("Getting student by ID",slog.String("userId",id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err!= nil {
            response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
            return
        }
        student, err := storage.GetStudentById(intId)
		
        if err!= nil {
			slog.Error("Error getting user",slog.String("id",id))
            response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
            return
        }

		response.WriteJson(w, http.StatusOK, student)

	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		slog.Info("Getting list of students")

        students, err := storage.GetStudents()		
        if err!= nil {
            response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
            return
        }
		response.WriteJson(w, http.StatusOK, students)
	}	
}
