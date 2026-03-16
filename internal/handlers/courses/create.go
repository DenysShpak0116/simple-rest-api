package courses

import (
	"log/slog"
	"net/http"
	"simple-rest-api/internal/handlers/responses"
	"simple-rest-api/internal/handlers/types"
	"simple-rest-api/internal/models"
	"simple-rest-api/pkg/http/render"
)

type CourseCreator interface {
	CreateCourse(course *models.Course) (int, error)
}

type TeacherValidator interface {
	GetTeacherByID(id int) (*models.Teacher, error)
}

func Create(logger *slog.Logger, creator CourseCreator, teacherValidator TeacherValidator) http.Handler {
	log := logger.With(slog.String("operation", "httpserver.handleCourseCreate"))

	type createCourseRequest struct {
		Title       string `json:"title" validate:"required,min=2,max=50"`
		Description string `json:"description" validate:"min=10,max=50"`
		TeacherID   int    `json:"teacher_id" validate:"required,gte=1"`
	}

	type courseResponse struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		TeacherID   *int   `json:"teacher_id"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, problems, err := render.DecodeValid[createCourseRequest](r)

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

		log.Debug("request validated", slog.String("title", req.Title), slog.Int("teacher_id", req.TeacherID))

		_, err = teacherValidator.GetTeacherByID(req.TeacherID)
		if err != nil {
			log.Warn("teacher not found", slog.Int("teacher_id", req.TeacherID))
			responses.Error(w, r, http.StatusBadRequest, "teacher not found")
			return
		}

		course := &models.Course{
			Title:       req.Title,
			Description: req.Description,
			TeacherID:   &req.TeacherID,
		}

		id, err := creator.CreateCourse(course)
		if err != nil {
			log.Error("failed to create course", slog.Any("error", err))
			responses.Error(w, r, http.StatusInternalServerError, "failed to create course")
			return
		}

		course.ID = id
		log.Info("course created successfully", slog.Int("id", id))

		render.Encode(w, r, http.StatusCreated, types.ApiResponse[courseResponse]{
			Message: "course created successfully",
			Data: &courseResponse{
				ID:          course.ID,
				Title:       course.Title,
				Description: course.Description,
				TeacherID:   course.TeacherID,
			},
		})
	})
}
