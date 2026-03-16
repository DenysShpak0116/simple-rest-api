package teachers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"simple-rest-api/internal/handlers/responses"
	"simple-rest-api/internal/repository"
)

type TeacherDeleter interface {
	DeleteTeacher(id int) error
}

func Delete(logger *slog.Logger, deleter TeacherDeleter) http.Handler {
	log := logger.With(slog.String("operation", "httpserver.handleTeacherDelete"))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			log.Warn("invalid teacher id", slog.String("id", idStr))
			responses.Error(w, r, http.StatusBadRequest, "invalid teacher id")
			return
		}

		log.Debug("deleting teacher", slog.Int("id", id))

		err = deleter.DeleteTeacher(id)
		if err != nil {
			if errors.Is(err, repository.ErrTeacherNotFound) {
				log.Warn("teacher not found", slog.Int("id", id))
				responses.Error(w, r, http.StatusNotFound, "teacher not found")
				return
			}
			log.Error("failed to delete teacher", slog.String("error", err.Error()))
			responses.Error(w, r, http.StatusInternalServerError, "failed to delete teacher")
			return
		}

		log.Info("teacher deleted successfully", slog.Int("id", id))

		w.WriteHeader(http.StatusNoContent)
	})
}
