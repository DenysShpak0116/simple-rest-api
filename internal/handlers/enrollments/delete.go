package enrollments

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"simple-rest-api/internal/handlers/responses"
	"simple-rest-api/internal/repository"
)

type EnrollmentDeleter interface {
	DeleteEnrollment(studentID, courseID int) error
}

func Delete(logger *slog.Logger, deleter EnrollmentDeleter, studentValidator StudentValidator, courseValidator CourseValidator) http.Handler {
	log := logger.With(slog.String("operation", "httpserver.handleEnrollmentDelete"))

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

		err = deleter.DeleteEnrollment(studentID, courseID)
		if err != nil {
			if errors.Is(err, repository.ErrEnrollmentNotFound) {
				log.Warn("enrollment not found", slog.Int("student_id", studentID), slog.Int("course_id", courseID))
				responses.Error(w, r, http.StatusNotFound, "enrollment not found")
				return
			}
			log.Error("failed to delete enrollment", slog.Any("error", err))
			responses.Error(w, r, http.StatusInternalServerError, "failed to delete enrollment")
			return
		}

		log.Info("enrollment deleted successfully", slog.Int("student_id", studentID), slog.Int("course_id", courseID))

		w.WriteHeader(http.StatusNoContent)
	})
}
