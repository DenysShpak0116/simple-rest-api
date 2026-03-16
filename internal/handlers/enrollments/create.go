package enrollments

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

type EnrollmentCreator interface {
	CreateEnrollment(enrollment *models.Enrollment) error
}

type StudentValidator interface {
	GetStudentByID(id int) (*models.Student, error)
}

type CourseValidator interface {
	GetCourseByID(id int) (*models.Course, error)
}

func Create(logger *slog.Logger, creator EnrollmentCreator, studentValidator StudentValidator, courseValidator CourseValidator) http.Handler {
	log := logger.With(slog.String("operation", "httpserver.handleEnrollmentCreate"))

	type createEnrollmentRequest struct {
		StudentID int `json:"student_id" validate:"required,gte=1"`
		CourseID  int `json:"course_id" validate:"required,gte=1"`
	}

	type enrollmentResponse struct {
		StudentID int `json:"student_id"`
		CourseID  int `json:"course_id"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		studentIDStr := r.PathValue("id")
		courseIDStr := r.PathValue("course_id")

		studentID, err := strconv.Atoi(studentIDStr)
		if err != nil || studentID <= 0 {
			log.Warn("invalid student id", slog.String("id", studentIDStr))
			responses.Error(w, r, http.StatusBadRequest, "invalid student id")
			return
		}

		courseID, err := strconv.Atoi(courseIDStr)
		if err != nil || courseID <= 0 {
			log.Warn("invalid course id", slog.String("course_id", courseIDStr))
			responses.Error(w, r, http.StatusBadRequest, "invalid course id")
			return
		}

		log.Debug("path parameters validated", slog.Int("student_id", studentID), slog.Int("course_id", courseID))

		_, err = studentValidator.GetStudentByID(studentID)
		if err != nil {
			log.Warn("student not found", slog.Int("student_id", studentID))
			responses.Error(w, r, http.StatusBadRequest, "student not found")
			return
		}

		_, err = courseValidator.GetCourseByID(courseID)
		if err != nil {
			log.Warn("course not found", slog.Int("course_id", courseID))
			responses.Error(w, r, http.StatusBadRequest, "course not found")
			return
		}

		enrollment := &models.Enrollment{
			StudentID: studentID,
			CourseID:  courseID,
		}

		err = creator.CreateEnrollment(enrollment)
		if err != nil {
			if errors.Is(err, repository.ErrEnrollmentConflict) {
				log.Warn("student already enrolled", slog.Int("student_id", studentID), slog.Int("course_id", courseID))
				responses.Error(w, r, http.StatusConflict, "student is already enrolled in this course")
				return
			}
			log.Error("failed to create enrollment", slog.Any("error", err))
			responses.Error(w, r, http.StatusInternalServerError, "failed to create enrollment")
			return
		}

		log.Info("enrollment created successfully", slog.Int("student_id", studentID), slog.Int("course_id", courseID))

		render.Encode(w, r, http.StatusCreated, types.ApiResponse[enrollmentResponse]{
			Message: "enrollment created successfully",
			Data: &enrollmentResponse{
				StudentID: enrollment.StudentID,
				CourseID:  enrollment.CourseID,
			},
		})
	})
}
