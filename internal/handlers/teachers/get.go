package teachers

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

type TeacherGetter interface {
	GetTeacherByID(id int) (*models.Teacher, error)
}

func Get(logger *slog.Logger, getter TeacherGetter) http.Handler {
	log := logger.With(slog.String("operation", "httpserver.handleTeacherGetByID"))

	type teacherResponse struct {
		ID         int    `json:"id"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Department string `json:"department"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			log.Warn("invalid teacher id", slog.String("id", idStr))
			responses.Error(w, r, http.StatusBadRequest, "invalid teacher id")
			return
		}

		log.Debug("fetching teacher", slog.Int("id", id))

		teacher, err := getter.GetTeacherByID(id)
		if err != nil {
			if errors.Is(err, repository.ErrTeacherNotFound) {
				log.Warn("teacher not found", slog.Int("id", id))
				responses.Error(w, r, http.StatusNotFound, "teacher not found")
				return
			}
			log.Error("failed to get teacher", slog.String("error", err.Error()))
			responses.Error(w, r, http.StatusInternalServerError, "failed to retrieve teacher")
			return
		}

		log.Info("teacher retrieved successfully", slog.Int("id", id))

		render.Encode(w, r, http.StatusOK, types.ApiResponse[teacherResponse]{
			Message: "teacher retrieved successfully",
			Data: &teacherResponse{
				ID:         teacher.ID,
				FirstName:  teacher.FirstName,
				LastName:   teacher.LastName,
				Department: teacher.Department,
			},
		})
	})
}
