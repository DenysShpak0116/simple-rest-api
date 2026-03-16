package courses

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

type CourseGetter interface {
	GetCourseByID(id int) (*models.Course, error)
}

func Get(logger *slog.Logger, getter CourseGetter) http.Handler {
	log := logger.With(slog.String("operation", "httpserver.handleCourseGetByID"))

	type courseResponse struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		TeacherID   *int   `json:"teacher_id"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			log.Warn("invalid course id", slog.String("id", idStr))
			responses.Error(w, r, http.StatusBadRequest, "invalid course id")
			return
		}

		log.Debug("fetching course", slog.Int("id", id))

		course, err := getter.GetCourseByID(id)
		if err != nil {
			if errors.Is(err, repository.ErrCourseNotFound) {
				log.Warn("course not found", slog.Int("id", id))
				responses.Error(w, r, http.StatusNotFound, "course not found")
				return
			}
			log.Error("failed to get course", slog.String("error", err.Error()))
			responses.Error(w, r, http.StatusInternalServerError, "failed to retrieve course")
			return
		}

		log.Info("course retrieved successfully", slog.Int("id", id))

		render.Encode(w, r, http.StatusOK, types.ApiResponse[courseResponse]{
			Message: "course retrieved successfully",
			Data: &courseResponse{
				ID:          course.ID,
				Title:       course.Title,
				Description: course.Description,
				TeacherID:   course.TeacherID,
			},
		})
	})
}
