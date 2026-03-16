package courses

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"simple-rest-api/internal/handlers/responses"
	"simple-rest-api/internal/repository"
)

type CourseDeleter interface {
	DeleteCourse(id int) error
}

func Delete(logger *slog.Logger, deleter CourseDeleter) http.Handler {
	log := logger.With(slog.String("operation", "httpserver.handleCourseDelete"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			log.Warn("invalid course id", slog.String("id", idStr))
			responses.Error(w, r, http.StatusBadRequest, "invalid course id")
			return
		}

		log.Debug("deleting course", slog.Int("id", id))

		err = deleter.DeleteCourse(id)
		if err != nil {
			if errors.Is(err, repository.ErrCourseNotFound) {
				log.Warn("course not found", slog.Int("id", id))
				responses.Error(w, r, http.StatusNotFound, "course not found")
				return
			}
			log.Error("failed to delete course", slog.String("error", err.Error()))
			responses.Error(w, r, http.StatusInternalServerError, "failed to delete course")
			return
		}

		log.Info("course deleted successfully", slog.Int("id", id))

		w.WriteHeader(http.StatusNoContent)
	})
}
