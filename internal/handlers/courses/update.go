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

type CourseUpdater interface {
	UpdateCourse(course *models.Course) error
}

type CourseValidator interface {
	GetCourseByID(id int) (*models.Course, error)
}

func Update(logger *slog.Logger, updater CourseUpdater, courseValidator CourseValidator, teacherValidator TeacherValidator) http.Handler {
	log := logger.With(slog.String("operation", "httpserver.handleCourseUpdate"))

	type updateCourseRequest struct {
		Title       string `json:"title" validate:"required,min=2,max=50"`
		Description string `json:"description" validate:"min=10,max=50"`
		TeacherID   *int   `json:"teacher_id" validate:"omitempty,gte=1"`
	}

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

		_, err = courseValidator.GetCourseByID(id)
		if err != nil {
			if errors.Is(err, repository.ErrCourseNotFound) {
				log.Warn("course not found", slog.Int("id", id))
				responses.Error(w, r, http.StatusNotFound, "course not found")
				return
			}
			log.Error("failed to validate course", slog.Any("error", err))
			responses.Error(w, r, http.StatusInternalServerError, "failed to validate course")
			return
		}

		req, problems, err := render.DecodeValid[updateCourseRequest](r)
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

		if req.TeacherID != nil {
			_, err = teacherValidator.GetTeacherByID(*req.TeacherID)
			if err != nil {
				if errors.Is(err, repository.ErrTeacherNotFound) {
					log.Warn("teacher not found", slog.Int("teacher_id", *req.TeacherID))
					responses.Error(w, r, http.StatusBadRequest, "teacher not found")
					return
				}
				log.Error("failed to validate teacher", slog.Any("error", err))
				responses.Error(w, r, http.StatusInternalServerError, "failed to validate teacher")
				return
			}
		}

		log.Debug("updating course", slog.Int("id", id))

		course := &models.Course{
			ID:          id,
			Title:       req.Title,
			Description: req.Description,
			TeacherID:   req.TeacherID,
		}

		err = updater.UpdateCourse(course)
		if err != nil {
			if errors.Is(err, repository.ErrCourseNotFound) {
				log.Warn("course not found", slog.Int("id", id))
				responses.Error(w, r, http.StatusNotFound, "course not found")
				return
			}
			log.Error("failed to update course", slog.String("error", err.Error()))
			responses.Error(w, r, http.StatusInternalServerError, "failed to update course")
			return
		}

		log.Info("course updated successfully", slog.Int("id", id))

		render.Encode(w, r, http.StatusOK, types.ApiResponse[courseResponse]{
			Message: "course updated successfully",
			Data: &courseResponse{
				ID:          course.ID,
				Title:       course.Title,
				Description: course.Description,
				TeacherID:   course.TeacherID,
			},
		})
	})
}
