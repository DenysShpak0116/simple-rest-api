package courses

import (
	"log/slog"
	"net/http"
	"simple-rest-api/internal/handlers/responses"
	"simple-rest-api/internal/handlers/types"
	"simple-rest-api/internal/models"
	"simple-rest-api/pkg/http/render"
)

type CoursesLister interface {
	GetAllCourses() ([]*models.Course, error)
}

func List(logger *slog.Logger, lister CoursesLister) http.Handler {
	log := logger.With(slog.String("operation", "httpserver.handleCoursesGet"))

	type courseResponse struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		TeacherID   *int   `json:"teacher_id"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		courses, err := lister.GetAllCourses()
		if err != nil {
			log.Error("failed to get courses", slog.String("error", err.Error()))
			responses.Error(w, r, http.StatusInternalServerError, "failed to retrieve courses")
			return
		}

		log.Info("courses retrieved successfully", slog.Int("count", len(courses)))

		courseResponses := make([]courseResponse, len(courses))
		for i, course := range courses {
			courseResponses[i] = courseResponse{
				ID:          course.ID,
				Title:       course.Title,
				Description: course.Description,
				TeacherID:   course.TeacherID,
			}
		}

		render.Encode(w, r, http.StatusOK, types.ApiResponse[[]courseResponse]{
			Message: "courses retrieved successfully",
			Data:    &courseResponses,
		})
	})
}
