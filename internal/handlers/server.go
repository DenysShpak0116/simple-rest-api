package httpserver

import (
	"log/slog"
	"net/http"
	"simple-rest-api/pkg/http/middleware"
	loggerMw "simple-rest-api/pkg/http/middleware/logger"

	_ "simple-rest-api/docs"
)

func NewServer(
	logger *slog.Logger,
	studentRepo StudentRepository,
	teacherRepo TeacherRepository,
	courseRepo CourseRepository,
	enrollmentRepo EnrollmentRepository,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		logger,
		studentRepo,
		teacherRepo,
		courseRepo,
		enrollmentRepo,
	)
	var handler http.Handler = mux

	handler = middleware.Chain(handler, loggerMw.New(logger))

	return handler
}
