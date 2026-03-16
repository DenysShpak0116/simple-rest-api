package students

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"simple-rest-api/internal/handlers/responses"
	"simple-rest-api/internal/handlers/types"
	"simple-rest-api/internal/models"
	"simple-rest-api/internal/repository"
	"simple-rest-api/pkg/http/render"
)

type StudentGetter interface {
	GetStudentByID(id int) (*models.Student, error)
}

func Get(logger *slog.Logger, getter StudentGetter) http.Handler {
	log := logger.With(slog.String("operation", "httpserver.handleStudentGetByID"))

	type studentResponse struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			log.Warn("invalid student id", slog.String("id", idStr))
			responses.Error(w, r, http.StatusBadRequest, "invalid student id")
			return
		}

		log.Debug("fetching student", slog.Int("id", id))

		student, err := getter.GetStudentByID(id)
		if err != nil {
			if errors.Is(err, repository.ErrStudentNotFound) {
				log.Warn("student not found", slog.Int("id", id))
				responses.Error(w, r, http.StatusNotFound, "student not found")
				return
			}
			log.Error("failed to get student", slog.String("error", err.Error()))
			responses.Error(w, r, http.StatusInternalServerError, "failed to retrieve student")
			return
		}

		log.Info("student retrieved successfully", slog.Int("id", id))

		render.Encode(w, r, http.StatusOK, types.ApiResponse[studentResponse]{
			Message: "student retrieved successfully",
			Data: &studentResponse{
				ID:        student.ID,
				FirstName: student.FirstName,
				LastName:  student.LastName,
				Email:     student.Email,
			},
		})
	})
}
