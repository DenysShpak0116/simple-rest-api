package teachers

import (
	"log/slog"
	"net/http"

	"simple-rest-api/internal/handlers/responses"
	"simple-rest-api/internal/handlers/types"
	"simple-rest-api/internal/models"
	"simple-rest-api/pkg/http/render"
)

type TeachersLister interface {
	GetAllTeachers() ([]*models.Teacher, error)
}

func List(logger *slog.Logger, lister TeachersLister) http.Handler {
	log := logger.With(slog.String("operation", "httpserver.handleTeachersGet"))

	type TeacherResponse struct {
		ID         int    `json:"id" example:"1"`
		FirstName  string `json:"first_name" example:"John"`
		LastName   string `json:"last_name" example:"Smith"`
		Department string `json:"department" example:"Mathematics"`
	}

	type teacherResponse struct {
		ID         int    `json:"id"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Department string `json:"department"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		teachers, err := lister.GetAllTeachers()
		if err != nil {
			log.Error("failed to get teachers", slog.String("error", err.Error()))
			responses.Error(w, r, http.StatusInternalServerError, "failed to retrieve teachers")
			return
		}

		log.Info("teachers retrieved successfully", slog.Int("count", len(teachers)))

		teacherResponses := make([]teacherResponse, len(teachers))
		for i, teacher := range teachers {
			teacherResponses[i] = teacherResponse{
				ID:         teacher.ID,
				FirstName:  teacher.FirstName,
				LastName:   teacher.LastName,
				Department: teacher.Department,
			}
		}

		render.Encode(w, r, http.StatusOK, types.ApiResponse[[]teacherResponse]{
			Message: "teachers retrieved successfully",
			Data:    &teacherResponses,
		})
	})
}
