package students

import (
	"log/slog"
	"net/http"

	"simple-rest-api/internal/handlers/responses"
	"simple-rest-api/internal/handlers/types"
	"simple-rest-api/internal/models"
	"simple-rest-api/pkg/http/render"
)

type StudentsLister interface {
	GetAllStudents() ([]*models.Student, error)
}

func List(logger *slog.Logger, lister StudentsLister) http.Handler {
	log := logger.With(slog.String("operation", "httpserver.handleStudentsGet"))

	type StudentResponse struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	type studentResponse struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		students, err := lister.GetAllStudents()
		if err != nil {
			log.Error("failed to get students", slog.String("error", err.Error()))
			responses.Error(w, r, http.StatusInternalServerError, "failed to retrieve students")
			return
		}

		log.Info("students retrieved successfully", slog.Int("count", len(students)))

		studentResponses := make([]studentResponse, len(students))
		for i, student := range students {
			studentResponses[i] = studentResponse{
				ID:        student.ID,
				FirstName: student.FirstName,
				LastName:  student.LastName,
				Email:     student.Email,
			}
		}

		render.Encode(w, r, http.StatusOK, types.ApiResponse[[]studentResponse]{
			Message: "students retrieved successfully",
			Data:    &studentResponses,
		})
	})
}
