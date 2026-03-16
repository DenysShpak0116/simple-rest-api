package students

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"simple-rest-api/internal/handlers/responses"
	"simple-rest-api/internal/repository"
)

type StudentDeleter interface {
	DeleteStudent(id int) error
}

func Delete(logger *slog.Logger, deleter StudentDeleter) http.Handler {
	log := logger.With(slog.String("operation", "httpserver.handleStudentDelete"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			log.Warn("invalid student id", slog.String("id", idStr))
			responses.Error(w, r, http.StatusBadRequest, "invalid student id")
			return
		}

		log.Debug("deleting student", slog.Int("id", id))

		err = deleter.DeleteStudent(id)
		if err != nil {
			if errors.Is(err, repository.ErrStudentNotFound) {
				log.Warn("student not found", slog.Int("id", id))
				responses.Error(w, r, http.StatusNotFound, "student not found")
				return
			}
			log.Error("failed to delete student", slog.String("error", err.Error()))
			responses.Error(w, r, http.StatusInternalServerError, "failed to delete student")
			return
		}

		log.Info("student deleted successfully", slog.Int("id", id))

		w.WriteHeader(http.StatusNoContent)
	})
}
