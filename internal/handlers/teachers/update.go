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

type TeacherUpdater interface {
	UpdateTeacher(teacher *models.Teacher) error
}

func Update(logger *slog.Logger, updater TeacherUpdater) http.Handler {
	log := logger.With(slog.String("operation", "httpserver.handleTeacherUpdate"))

	type updateTeacherRequest struct {
		FirstName  string `json:"first_name" validate:"required,min=2,max=20"`
		LastName   string `json:"last_name" validate:"required,min=2,max=20"`
		Department string `json:"department" validate:"required"`
	}

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

		req, problems, err := render.DecodeValid[updateTeacherRequest](r)
		if len(problems) > 0 {
			log.Warn("validation failed", slog.Any("problems", problems))
			responses.ValidationError(w, r, problems)
			return
		}

		if err != nil {
			log.Error("decode request failed", slog.Any("error", err))
			responses.Error(w, r, http.StatusBadRequest, "invalid request data")
			return
		}

		log.Debug("updating teacher", slog.Int("id", id))

		teacher := &models.Teacher{
			ID:         id,
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			Department: req.Department,
		}

		err = updater.UpdateTeacher(teacher)
		if err != nil {
			if errors.Is(err, repository.ErrTeacherNotFound) {
				log.Warn("teacher not found", slog.Int("id", id))
				responses.Error(w, r, http.StatusNotFound, "teacher not found")
				return
			}
			log.Error("failed to update teacher", slog.String("error", err.Error()))
			responses.Error(w, r, http.StatusInternalServerError, "failed to update teacher")
			return
		}

		log.Info("teacher updated successfully", slog.Int("id", id))

		render.Encode(w, r, http.StatusOK, types.ApiResponse[teacherResponse]{
			Message: "teacher updated successfully",
			Data: &teacherResponse{
				ID:         teacher.ID,
				FirstName:  teacher.FirstName,
				LastName:   teacher.LastName,
				Department: teacher.Department,
			},
		})
	})
}
