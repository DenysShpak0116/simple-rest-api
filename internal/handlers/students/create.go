package students

import (
	"errors"
	"log/slog"
	"net/http"

	"simple-rest-api/internal/handlers/responses"
	"simple-rest-api/internal/handlers/types"
	"simple-rest-api/internal/models"
	"simple-rest-api/internal/repository"
	"simple-rest-api/pkg/http/render"
)

type StudentCreator interface {
	CreateStudent(student *models.Student) (int, error)
}

func Create(logger *slog.Logger, creator StudentCreator) http.Handler {
	log := logger.With(slog.String("operation", "httpserver.handleStudentCreate"))

	type createStudentRequest struct {
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
		req, problems, err := render.DecodeValid[createStudentRequest](r)
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

		log.Debug("request validated", slog.String("email", req.Email))

		student := &models.Student{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
		}

		id, err := creator.CreateStudent(student)
		if err != nil {
			if errors.Is(err, repository.ErrDuplicateEmail) {
				log.Warn("email already exists", slog.String("email", req.Email))
				responses.Error(w, r, http.StatusConflict, "email already exists")
				return
			}
			log.Error("failed to create student", slog.Any("error", err))
			responses.Error(w, r, http.StatusInternalServerError, "failed to create student")
			return
		}

		student.ID = id
		log.Info("student created successfully", slog.Int("id", id))

		render.Encode(w, r, http.StatusCreated, types.ApiResponse[studentResponse]{
			Message: "student created successfully",
			Data: &studentResponse{
				ID:        student.ID,
				FirstName: student.FirstName,
				LastName:  student.LastName,
				Email:     student.Email,
			},
		})
	})
}
