package teachers

import (
	"log/slog"
	"net/http"

	"simple-rest-api/internal/handlers/responses"
	"simple-rest-api/internal/handlers/types"
	"simple-rest-api/internal/models"
	"simple-rest-api/pkg/http/render"
)

type TeacherCreator interface {
	CreateTeacher(teacher *models.Teacher) (int, error)
}

func Create(logger *slog.Logger, creator TeacherCreator) http.Handler {
	log := logger.With(slog.String("operation", "httpserver.handleTeacherCreate"))

	type createTeacherRequest struct {
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
		req, problems, err := render.DecodeValid[createTeacherRequest](r)
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

		log.Debug("request validated", slog.String("department", req.Department))

		teacher := &models.Teacher{
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			Department: req.Department,
		}

		id, err := creator.CreateTeacher(teacher)
		if err != nil {
			log.Error("failed to create teacher", slog.Any("error", err))
			responses.Error(w, r, http.StatusInternalServerError, "failed to create teacher")
			return
		}

		teacher.ID = id
		log.Info("teacher created successfully", slog.Int("id", id))

		render.Encode(w, r, http.StatusCreated, types.ApiResponse[teacherResponse]{
			Message: "teacher created successfully",
			Data: &teacherResponse{
				ID:         teacher.ID,
				FirstName:  teacher.FirstName,
				LastName:   teacher.LastName,
				Department: teacher.Department,
			},
		})
	})
}
