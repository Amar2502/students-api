package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/amar2502/students-api/internal/storage"
	"github.com/amar2502/students-api/internal/types"
	"github.com/amar2502/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10" 
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		slog.Info("Decoded the request body", "student", student)

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
 
		slog.Info("Creating a new student")


		// Validate the request body
		if err := validator.New().Struct(student); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)

		slog.Info("User created sucessfully", slog.String("LastId", fmt.Sprint(lastId)))

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"lastId": lastId})
	}
}

func GetbyId(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting student by id", slog.String("id", id))

		intId, err := 	strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Info("error converting id into int64")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById(intId)
		if err != nil {
			slog.Info("error getting student by id")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return	
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Getting students")

		students, err := storage.GetStudent()
		if err!=nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		slog.Info("Students fetched")

		if students == nil {
			response.WriteJson(w, http.StatusInternalServerError, students)
			return
		}

		slog.Info("student not null")

		response.WriteJson(w, http.StatusOK, students)

	}
}