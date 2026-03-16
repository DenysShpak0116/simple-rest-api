package httpserver

import (
	"log/slog"
	"net/http"
	"simple-rest-api/internal/handlers/courses"
	"simple-rest-api/internal/handlers/enrollments"
	"simple-rest-api/internal/handlers/students"
	"simple-rest-api/internal/handlers/teachers"

	httpSwagger "github.com/swaggo/http-swagger"
)

func addRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	studentRepo StudentRepository,
	teacherRepo TeacherRepository,
	courseRepo CourseRepository,
	enrollmentRepo EnrollmentRepository,
) {
	mux.Handle("/", http.NotFoundHandler())

	mux.Handle("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8081/swagger/doc.json"),
		httpSwagger.DocExpansion("list"),
	))
	mux.Handle("GET /api/docs", http.RedirectHandler("/swagger/", http.StatusMovedPermanently))

	mux.Handle("GET /students", students.List(logger, studentRepo))
	mux.Handle("GET /students/{id}", students.Get(logger, studentRepo))
	mux.Handle("POST /students", students.Create(logger, studentRepo))
	mux.Handle("PUT /students/{id}", students.Update(logger, studentRepo))
	mux.Handle("DELETE /students/{id}", students.Delete(logger, studentRepo))

	mux.Handle("GET /teachers", teachers.List(logger, teacherRepo))
	mux.Handle("GET /teachers/{id}", teachers.Get(logger, teacherRepo))
	mux.Handle("POST /teachers", teachers.Create(logger, teacherRepo))
	mux.Handle("PUT /teachers/{id}", teachers.Update(logger, teacherRepo))
	mux.Handle("DELETE /teachers/{id}", teachers.Delete(logger, teacherRepo))

	mux.Handle("GET /courses", courses.List(logger, courseRepo))
	mux.Handle("GET /courses/{id}", courses.Get(logger, courseRepo))
	mux.Handle("POST /courses", courses.Create(logger, courseRepo, teacherRepo))
	mux.Handle("PUT /courses/{id}", courses.Update(logger, courseRepo, courseRepo, teacherRepo))
	mux.Handle("DELETE /courses/{id}", courses.Delete(logger, courseRepo))

	mux.Handle("POST /students/{id}/courses/{course_id}", enrollments.Create(logger, enrollmentRepo, studentRepo, courseRepo))
	mux.Handle("DELETE /students/{id}/courses/{course_id}", enrollments.Delete(logger, enrollmentRepo, studentRepo, courseRepo))
}
