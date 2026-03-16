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

type StudentUpdater interface {
	UpdateStudent(student *models.Student) error
}

func Update(logger *slog.Logger, updater StudentUpdater) http.Handler {
	log := logger.With(slog.String("operation", "httpserver.handleStudentUpdate"))

	type updateStudentRequest struct {
		FirstName string `json:"first_name" validate:"required,min=2,max=20"`
		LastName  string `json:"last_name" validate:"required,min=2,max=20"`
		Email     string `json:"email" validate:"required,email"`
	}

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

		req, problems, err := render.DecodeValid[updateStudentRequest](r)
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

		log.Debug("updating student", slog.Int("id", id))

		student := &models.Student{
			ID:        id,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
		}

		err = updater.UpdateStudent(student)
		if err != nil {
			if errors.Is(err, repository.ErrStudentNotFound) {
				log.Warn("student not found", slog.Int("id", id))
				responses.Error(w, r, http.StatusNotFound, "student not found")
				return
			}
			log.Error("failed to update student", slog.String("error", err.Error()))
			responses.Error(w, r, http.StatusInternalServerError, "failed to update student")
			return
		}

		log.Info("student updated successfully", slog.Int("id", id))

		render.Encode(w, r, http.StatusOK, types.ApiResponse[studentResponse]{
			Message: "student updated successfully",
			Data: &studentResponse{
				ID:        student.ID,
				FirstName: student.FirstName,
				LastName:  student.LastName,
				Email:     student.Email,
			},
		})
	})
}
